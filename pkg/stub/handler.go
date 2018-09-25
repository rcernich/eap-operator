package stub

import (
	"strconv"
	"context"
	"fmt"
	"strings"

	"github.com/jboss-openshift/eap-operator/pkg/apis/eap/v1alpha1"
	appsv1 "github.com/openshift/api/apps/v1"
	buildv1 "github.com/openshift/api/build/v1"
	imagev1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.EAPApplicationConfig:
		// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
		// All secondary resources must have the CR set as their OwnerReference for this to be the case
		if event.Deleted {
			return nil
		}

		owner := asOwner(o)
		labels := createLabels(o)

		// output image
		err := syncOutputImage(o, &labels, &owner)
		if err != nil {
			logrus.Errorf("Failed to update output ImageStream for %v:%v: %v", o.GetNamespace(), o.GetName(), err)
			return err
		}

		// build configuration
		err = syncBuildConfig(o, &labels, &owner)
		if err != nil {
			logrus.Errorf("Failed to update BuildConfig for %v:%v: %v", o.GetNamespace(), o.GetName(), err)
			return err
		}

		// deployment configuration
		err = syncDeploymentConfig(o, &labels, &owner)
		if err != nil {
			logrus.Errorf("Failed to update DeploymentConfig for %v:%v: %v", o.GetNamespace(), o.GetName(), err)
			return err
		}

		// services
		err = syncServices(o, &labels, &owner)
		if err != nil {
			logrus.Errorf("Failed to update Service for %v:%v: %v", o.GetNamespace(), o.GetName(), err)
			return err
		}

		// routes
		err = syncRoutes(o, &labels, &owner)
		if err != nil {
			logrus.Errorf("Failed to update Service for %v:%v: %v", o.GetNamespace(), o.GetName(), err)
			return err
		}
	}
	return nil
}

// asOwner returns an OwnerReference set as the EAP Application CR
func asOwner(config *v1alpha1.EAPApplicationConfig) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: config.APIVersion,
		Kind:       config.Kind,
		Name:       config.Name,
		UID:        config.UID,
		Controller: &trueVar,
	}
}

func createLabels(config *v1alpha1.EAPApplicationConfig) map[string]string {
	return map[string]string{
		"app":      config.GetName(),
		"operator": "eap-operator",
	}
}

func outputImageName(objectMeta metav1.Object) string {
	return objectMeta.GetName()
}

func outputImageTag() string {
	return "latest"
}
func syncOutputImage(config *v1alpha1.EAPApplicationConfig, labels *map[string]string, owner *metav1.OwnerReference) error {
	imageStream := &imagev1.ImageStream{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ImageStream",
			APIVersion: imagev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
	}
	err := sdk.Get(imageStream)
	if err == nil {
		// nothing to do
		return nil
	} else if !errors.IsNotFound(err) {
		return fmt.Errorf("Error synchronizing ImageStream: %v", err)
	}
	imageStream.Labels = *labels
	imageStream.SetOwnerReferences(append(imageStream.GetOwnerReferences(), *owner))
	err = sdk.Create(imageStream)
	return err
}

func syncBuildConfig(config *v1alpha1.EAPApplicationConfig, labels *map[string]string, owner *metav1.OwnerReference) error {
	buildConfig := &buildv1.BuildConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildConfig",
			APIVersion: buildv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
	}

	var createOrUpdate func(object sdk.Object) error
	err := sdk.Get(buildConfig)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("Error synchronizing BuildConfig: %v", err)
		}
		createOrUpdate = sdk.Create
		buildConfig.SetOwnerReferences(append(buildConfig.GetOwnerReferences(), *owner))
		} else {
		// Update the existing BuildConfig
		createOrUpdate = sdk.Update
	}
	buildConfig.Labels = *labels
	switch config.Spec.BuildConfig.Type {
	case v1alpha1.MavenBuildType:
		processGenericBuildConfig(&config.Spec.BuildConfig.MavenBuildConfig.GenericBuildConfig, &buildConfig.Spec)
		processMavenOptions(config.Spec.BuildConfig.MavenBuildConfig.MavenOptions, &buildConfig.Spec)
	default:
		return fmt.Errorf("Unknown BuildType in config: %v", config.Spec.BuildConfig.Type)
	}
	buildConfig.Spec.Output.To = &corev1.ObjectReference{
		Kind: "ImageStreamTag",
		Name: outputImageName(config.GetObjectMeta()) + ":" + outputImageTag(),
	}
	if config.Spec.Image == nil {
		buildConfig.Spec.Strategy.SourceStrategy.From = corev1.ObjectReference{
			Kind: "ImageStreamTag",
			Name: "jboss-eap64-openshift:1.8",
			Namespace: "openshift",
		}
	} else {
		buildConfig.Spec.Strategy.SourceStrategy.From = *config.Spec.Image
	}
	err = createOrUpdate(buildConfig)
	return err
}

func processGenericBuildConfig(config *v1alpha1.GenericBuildConfig, buildConfig *buildv1.BuildConfigSpec) {
	var incrementalBuilds *bool
	if config.GenericOptions != nil && config.GenericOptions.IncrementalBuilds != nil {
		incrementalBuilds = config.GenericOptions.IncrementalBuilds
	}
	buildConfig.Source = config.BuildSource
	buildConfig.Strategy.Type = buildv1.SourceBuildStrategyType
	buildConfig.Strategy.SourceStrategy = &buildv1.SourceBuildStrategy{
		Env:         config.Env,
		ForcePull:   true,
		Incremental: incrementalBuilds,
	}
	processImageSourceMounts(config, buildConfig)
	processGenericOptions(config.GenericOptions, buildConfig)
}

func processImageSourceMounts(config *v1alpha1.GenericBuildConfig, buildConfig *buildv1.BuildConfigSpec) {
	if config.ImageSourceMountDir != nil && len(buildConfig.Source.Images) > 0 {
		contextDir := buildConfig.Source.ContextDir
		if len(contextDir) > 0 {
			contextDir += "/" + *config.ImageSourceMountDir
		} else {
			contextDir = *config.ImageSourceMountDir
		}
		if !strings.HasSuffix(contextDir, "/") {
			contextDir += "/"
		}
		for _, image := range buildConfig.Source.Images {
			for _, path := range image.Paths {
				path.DestinationDir = contextDir + path.DestinationDir
			}
		}
		imageSourceMountDir := *config.ImageSourceMountDir
		if !strings.HasSuffix(imageSourceMountDir, "/") {
			imageSourceMountDir += "/"
		}
		imageSourceMountDir += "*" // add in all the subdirectories
		buildConfig.Strategy.SourceStrategy.Env = append(buildConfig.Strategy.SourceStrategy.Env, corev1.EnvVar{
			Name:  "S2I_IMAGE_SOURCE_MOUNTS",
			Value: imageSourceMountDir,
		})
	}
}

func processGenericOptions(options *v1alpha1.GenericBuildOptions, buildConfig *buildv1.BuildConfigSpec) {
	if options == nil {
		return
	}
	envs := make([]corev1.EnvVar,0,8)
	if options.IncrementalBuilds != nil {
		envs = append(envs, corev1.EnvVar{
			Name: "S2I_ENABLE_INCREMENTAL_BUILDS",
			Value: strconv.FormatBool(*options.IncrementalBuilds),
		})
	}
	if options.ConfigurationPaths != nil {
		envs = append(envs, corev1.EnvVar{
			Name: "S2I_SOURCE_CONFIGURATION_DIR",
			Value: options.ConfigurationPaths.SourcePath,
		})
		envs = append(envs, corev1.EnvVar{
			Name: "S2I_TARGET_CONFIGURATION_DIR",
			Value: options.ConfigurationPaths.DestinationDir,
		})
	}
	if options.DataPaths != nil {
		envs = append(envs, corev1.EnvVar{
			Name: "S2I_SOURCE_DATA_DIR",
			Value: options.DataPaths.SourcePath,
		})
		envs = append(envs, corev1.EnvVar{
			Name: "S2I_TARGET_DATA_DIR",
			Value: options.DataPaths.DestinationDir,
		})
	}
	if options.DeploymentPaths != nil {
		envs = append(envs, corev1.EnvVar{
			Name: "S2I_SOURCE_DEPLOYMENTS_DIR",
			Value: options.DeploymentPaths.SourcePath,
		})
		envs = append(envs, corev1.EnvVar{
			Name: "S2I_TARGET_DEPLOYMENTS_DIR",
			Value: options.DeploymentPaths.DestinationDir,
		})
	}
	buildConfig.Strategy.SourceStrategy.Env = append(buildConfig.Strategy.SourceStrategy.Env, envs...)
}

func processMavenOptions(options *v1alpha1.MavenOptions, buildConfig *buildv1.BuildConfigSpec) {
	// TODO
}

func syncDeploymentConfig(config *v1alpha1.EAPApplicationConfig, labels *map[string]string, owner *metav1.OwnerReference) error {
	deployConfig := &appsv1.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DeploymentConfig",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
	}

	var createOrUpdate func(object sdk.Object) error
	err := sdk.Get(deployConfig)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("Error synchronizing DeploymentConfig: %v", err)
		}
		createOrUpdate = sdk.Create
		deployConfig.SetOwnerReferences(append(deployConfig.GetOwnerReferences(), *owner))
		} else {
		// Update the existing DeploymentConfig
		createOrUpdate = sdk.Update
	}
	deployConfig.Labels = *labels

	var timeout int64 = 60
	var memoryLimit resource.Quantity
	memoryLimit, err = resource.ParseQuantity("1Gi")
	if err != nil {
		return fmt.Errorf("Error parsing memory limit: %v", err)
	}
	deployConfig.Spec = appsv1.DeploymentConfigSpec{
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.DeploymentStrategyTypeRecreate,
		},
		Triggers: appsv1.DeploymentTriggerPolicies{
			appsv1.DeploymentTriggerPolicy{
				Type: "ImageChange",
				ImageChangeParams: &appsv1.DeploymentTriggerImageChangeParams{
					Automatic: true,
					ContainerNames: []string{
						config.GetName(),
					},
					From: corev1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: outputImageName(config) + ":" + outputImageTag(),
					},
				},
			},
		},
		Replicas: 1,
		Selector: *labels,
		Template: &corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name: config.GetName(),
				Labels: *labels,
			},
			Spec: corev1.PodSpec{
				TerminationGracePeriodSeconds: &timeout,
				Containers: []corev1.Container{
					corev1.Container{
						Name: config.GetName(),
						Image: outputImageName(config) + ":" + outputImageTag(),
						ImagePullPolicy: corev1.PullAlways,
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: memoryLimit,
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"/bin/bash",
										"-c",
										"/opt/eap/bin/readinessProbe.sh",
									},
								},
							},
						},
						Ports: []corev1.ContainerPort{
							corev1.ContainerPort{
								Name: "jolokia",
								ContainerPort: 8778,
								Protocol: corev1.ProtocolTCP,
							},
							corev1.ContainerPort{
								Name: "http",
								ContainerPort: 8080,
								Protocol: corev1.ProtocolTCP,
							},
						},
						Env: []corev1.EnvVar{
							corev1.EnvVar{
								Name: "JGROUPS_PING_PROTOCOL",
								Value: "openshift.DNS_PING",
							},
							corev1.EnvVar{
								Name: "OPENSHIFT_DNS_PING_SERVICE_NAME",
								Value: config.GetName() + "-ping",
							},
							corev1.EnvVar{
								Name: "OPENSHIFT_DNS_PING_SERVICE_PORT",
								Value: "8888",
							},
						},
					},
				},
			},
		},
	}
	err = createOrUpdate(deployConfig)
	return err
}

func syncServices(config *v1alpha1.EAPApplicationConfig, labels *map[string]string, owner *metav1.OwnerReference) error {
	appService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
	}

	clusterService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name + "-ping",
			Namespace: config.Namespace,
		},
	}

	var createOrUpdate func(object sdk.Object) error
	err := sdk.Get(appService)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("Error synchronizing application Service: %v", err)
		}
		createOrUpdate = sdk.Create
		appService.SetOwnerReferences(append(appService.GetOwnerReferences(), *owner))
		} else {
		// Update the existing Service
		createOrUpdate = sdk.Update
	}
	appService.Labels = *labels
	appService.Spec = corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			corev1.ServicePort{
				Name: "http",
				TargetPort: intstr.FromString("http"),
				Port: 8080,
			},
		},
		ClusterIP: appService.Spec.ClusterIP,
		Selector: *labels,
	}
	err = createOrUpdate(appService)
	if err != nil {
		return fmt.Errorf("Error synchronizing application Service: %v", err)
	}

	err = sdk.Get(clusterService)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("Error synchronizing cluster Service: %v", err)
		}
		createOrUpdate = sdk.Create
		clusterService.SetOwnerReferences(append(clusterService.GetOwnerReferences(), *owner))
		} else {
		// Update the existing Service
		createOrUpdate = sdk.Update
	}
	clusterService.Labels = *labels
	clusterService.Spec = corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			corev1.ServicePort{
				Name: "http",
				TargetPort: intstr.FromString("http"),
				Port: 8080,
			},
		},
		Selector: *labels,
		ClusterIP: corev1.ClusterIPNone,
		PublishNotReadyAddresses: true,
	}
	err = createOrUpdate(clusterService)
	if err != nil {
		return fmt.Errorf("Error synchronizing cluster Service: %v", err)
	}
	return nil
}

func syncRoutes(config *v1alpha1.EAPApplicationConfig, labels *map[string]string, owner *metav1.OwnerReference) error {
	routeConfig := &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: routev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
	}

	var createOrUpdate func(object sdk.Object) error
	err := sdk.Get(routeConfig)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("Error synchronizing Route: %v", err)
		}
		createOrUpdate = sdk.Create
		routeConfig.SetOwnerReferences(append(routeConfig.GetOwnerReferences(), *owner))
		} else {
		// Update the existing Route
		createOrUpdate = sdk.Update
	}
	routeConfig.Labels = *labels
	routeConfig.Spec = routev1.RouteSpec{
		To: routev1.RouteTargetReference{
			Kind: "Service",
			Name: config.GetName(),
		},
	}
	err = createOrUpdate(routeConfig)
	return err
}
