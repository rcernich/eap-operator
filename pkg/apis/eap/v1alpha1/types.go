package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	openshiftv1 "github.com/openshift/api/build/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EAPApplicationConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []EAPApplicationConfig `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EAPApplicationConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              EAPApplicationConfigSpec   `json:"spec"`
	Status            EAPApplicationConfigStatus `json:"status,omitempty"`
}

// EAPApplicationConfigSpec defines the configuration for an EAP based application
type EAPApplicationConfigSpec struct {
	BuildConfig ApplicationBuildConfig `json:"buildConfig"`
	// image is reference to an DockerImage, ImageStreamTag, or ImageStreamImage which is used instead of the default image
	Image *corev1.ObjectReference `json:"image,omitempty"`
}

// EAPApplicationConfigStatus defines the current status of the operator
type EAPApplicationConfigStatus struct {
	// Fill me
}

// ApplicationBuildConfig defines the configuration for building the user's application
type ApplicationBuildConfig struct {
	Type ApplicationBuildType `json:"type"`
	MavenBuildConfig *MavenBuildConfig `json:"mavenBuildConfig,omitempty"`
}

// ApplicationBuildType constants for various supported build types
type ApplicationBuildType string
const (
	// MavenBuildType instructs the operator to use the maven builder
	MavenBuildType ApplicationBuildType = "maven"
)

// GenericBuildConfig defines basic elements common to all build configuration types
type GenericBuildConfig struct {
	BuildSource openshiftv1.BuildSource `json:"source"`
	// Location within source directory where ImageSource entries should be mounted. This directory will be prefixed to
	// the destinationDir entries on any ImageSource specified in ApplicationBuildSource
	ImageSourceMountDir *string `json:"imageSourceMountDir,omitempty"`
	// Additional user specified environment variables
	Env []corev1.EnvVar `json:"env,omitempty"`
	GenericOptions *GenericBuildOptions `json:"genericOptions,omitempty"`
}

// GenericBuildOptions defines options used by all build types
type GenericBuildOptions struct {
	// specifies relative source path and absolute target path from which configration files are copied from/to
	ConfigurationPaths *openshiftv1.ImageSourcePath `json:"configurationPaths,omitempty"`
	// specifies relative source path and absolute target path from which data files are copied from/to
	DataPaths *openshiftv1.ImageSourcePath `json:"dataPaths,omitempty"`
	// specifies relative source path and absolute target path from which deployment files are copied from/to
	DeploymentPaths *openshiftv1.ImageSourcePath `json:"deploymentPaths,omitempty"`
	// specifies whether or not incremental builds should be enabled
	IncrementalBuilds *bool `json:"incrementalBuilds,omitempty"`
}

// MavenBuildConfig defines build configuration specific to maven builds
type MavenBuildConfig struct {
	GenericBuildConfig `json:",inline"`
	MavenOptions *MavenOptions `json:"mavenOptions,omitempty"`
}

// MavenOptions defines options to be used when running maven
type MavenOptions struct {
	// list of maven goals that should be executed.  defaults to image default,which is typically package.
	Goals []string `json:"goals,omitempty"`
	// specifies relative paths within source from which build artifacts are copied.  defaults to image default, which
	// is typically target.
	ArtifactPaths []string `json:"artifactPaths,omitempty"`
	// override default maven arguments
	MavenArgs *string `json:"mavenArgs,omitempty"`
	// arguments added to the default arguments
	AdditionalMavenArgs *string `json:"additionalMavenArgs,omitempty"`
	// clear local maven repository upon build completion.  this will be overridden by BuildOptions.IncrementalBuild.
	ClearRepo *bool `json:"clearRepo,omitempty"`
	// path to custom settings.xml file to use when building
	SettingsFile *string `json:"settingsFile,omitempty"`
	// path to local repository
	LocalRepositoryPath *string `json:"localRepositoryPath,omitempty"`
	// mirrors that should be configured for the build
	Mirrors []MavenMirror `json:"mirrors,omitempty"`
	// maven repositories that should be configured for the build
	Repositories []MavenRepository `json:"repositories,omitempty"`
}

// MavenMirror defines a maven mirror that should be added to settings.xml file
type MavenMirror struct {
	// the id to use for this mirror.  defaults to auto-generated unique id.
	ID *string `json:"id,omitempty"`
	// the URL of the mirror.
	URL string `json:"url"`
	// the repositories which are mirrored by this mirror.  defaults to external:*
	Of *string `json:"of,omitempty"`
}

// MavenRepository defines a maven repository that should be added to settings.xml file
type MavenRepository struct {
	// the id for this repository.  defaults to auto-generated unique id.
	ID *string `json:"id,omitempty"`
	// the URL for this repository.
	URL string `json:"url"`
	// the layout of the repository.  defaults to default.
	Layout *string `json:"layout,omitempty"`
	// are releases enabled?  defaults to true.
	ReleasesEnabled *bool `json:"releasesEnabled,omitempty"`
	// releases update policy.  defaults to always.
	ReleasesUpdatePolicy *string `json:"releasesUpdatePolicy,omitempty"`
	// are snapshots enabled? defaults to true.
	SnapshotsEnabled *bool `json:"snapshotsEnabled,omitempty"`
	// snapshots update policy.  defaults to always.
	SnapshotsUpdatePolicy *string `json:"snapshotsUpdatePolicy,omitempty"`
}
