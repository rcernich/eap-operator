// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	buildv1 "github.com/openshift/api/build/v1"
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationBuildConfig) DeepCopyInto(out *ApplicationBuildConfig) {
	*out = *in
	if in.MavenBuildConfig != nil {
		in, out := &in.MavenBuildConfig, &out.MavenBuildConfig
		*out = new(MavenBuildConfig)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationBuildConfig.
func (in *ApplicationBuildConfig) DeepCopy() *ApplicationBuildConfig {
	if in == nil {
		return nil
	}
	out := new(ApplicationBuildConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EAPApplicationConfig) DeepCopyInto(out *EAPApplicationConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EAPApplicationConfig.
func (in *EAPApplicationConfig) DeepCopy() *EAPApplicationConfig {
	if in == nil {
		return nil
	}
	out := new(EAPApplicationConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EAPApplicationConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EAPApplicationConfigList) DeepCopyInto(out *EAPApplicationConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EAPApplicationConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EAPApplicationConfigList.
func (in *EAPApplicationConfigList) DeepCopy() *EAPApplicationConfigList {
	if in == nil {
		return nil
	}
	out := new(EAPApplicationConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EAPApplicationConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EAPApplicationConfigSpec) DeepCopyInto(out *EAPApplicationConfigSpec) {
	*out = *in
	in.BuildConfig.DeepCopyInto(&out.BuildConfig)
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(v1.ObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EAPApplicationConfigSpec.
func (in *EAPApplicationConfigSpec) DeepCopy() *EAPApplicationConfigSpec {
	if in == nil {
		return nil
	}
	out := new(EAPApplicationConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EAPApplicationConfigStatus) DeepCopyInto(out *EAPApplicationConfigStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EAPApplicationConfigStatus.
func (in *EAPApplicationConfigStatus) DeepCopy() *EAPApplicationConfigStatus {
	if in == nil {
		return nil
	}
	out := new(EAPApplicationConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericBuildConfig) DeepCopyInto(out *GenericBuildConfig) {
	*out = *in
	in.BuildSource.DeepCopyInto(&out.BuildSource)
	if in.ImageSourceMountDir != nil {
		in, out := &in.ImageSourceMountDir, &out.ImageSourceMountDir
		*out = new(string)
		**out = **in
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.GenericOptions != nil {
		in, out := &in.GenericOptions, &out.GenericOptions
		*out = new(GenericBuildOptions)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericBuildConfig.
func (in *GenericBuildConfig) DeepCopy() *GenericBuildConfig {
	if in == nil {
		return nil
	}
	out := new(GenericBuildConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericBuildOptions) DeepCopyInto(out *GenericBuildOptions) {
	*out = *in
	if in.ConfigurationPaths != nil {
		in, out := &in.ConfigurationPaths, &out.ConfigurationPaths
		*out = new(buildv1.ImageSourcePath)
		**out = **in
	}
	if in.DataPaths != nil {
		in, out := &in.DataPaths, &out.DataPaths
		*out = new(buildv1.ImageSourcePath)
		**out = **in
	}
	if in.DeploymentPaths != nil {
		in, out := &in.DeploymentPaths, &out.DeploymentPaths
		*out = new(buildv1.ImageSourcePath)
		**out = **in
	}
	if in.IncrementalBuilds != nil {
		in, out := &in.IncrementalBuilds, &out.IncrementalBuilds
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericBuildOptions.
func (in *GenericBuildOptions) DeepCopy() *GenericBuildOptions {
	if in == nil {
		return nil
	}
	out := new(GenericBuildOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MavenBuildConfig) DeepCopyInto(out *MavenBuildConfig) {
	*out = *in
	in.GenericBuildConfig.DeepCopyInto(&out.GenericBuildConfig)
	if in.MavenOptions != nil {
		in, out := &in.MavenOptions, &out.MavenOptions
		*out = new(MavenOptions)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MavenBuildConfig.
func (in *MavenBuildConfig) DeepCopy() *MavenBuildConfig {
	if in == nil {
		return nil
	}
	out := new(MavenBuildConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MavenMirror) DeepCopyInto(out *MavenMirror) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Of != nil {
		in, out := &in.Of, &out.Of
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MavenMirror.
func (in *MavenMirror) DeepCopy() *MavenMirror {
	if in == nil {
		return nil
	}
	out := new(MavenMirror)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MavenOptions) DeepCopyInto(out *MavenOptions) {
	*out = *in
	if in.Goals != nil {
		in, out := &in.Goals, &out.Goals
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ArtifactPaths != nil {
		in, out := &in.ArtifactPaths, &out.ArtifactPaths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.MavenArgs != nil {
		in, out := &in.MavenArgs, &out.MavenArgs
		*out = new(string)
		**out = **in
	}
	if in.AdditionalMavenArgs != nil {
		in, out := &in.AdditionalMavenArgs, &out.AdditionalMavenArgs
		*out = new(string)
		**out = **in
	}
	if in.ClearRepo != nil {
		in, out := &in.ClearRepo, &out.ClearRepo
		*out = new(bool)
		**out = **in
	}
	if in.SettingsFile != nil {
		in, out := &in.SettingsFile, &out.SettingsFile
		*out = new(string)
		**out = **in
	}
	if in.LocalRepositoryPath != nil {
		in, out := &in.LocalRepositoryPath, &out.LocalRepositoryPath
		*out = new(string)
		**out = **in
	}
	if in.Mirrors != nil {
		in, out := &in.Mirrors, &out.Mirrors
		*out = make([]MavenMirror, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Repositories != nil {
		in, out := &in.Repositories, &out.Repositories
		*out = make([]MavenRepository, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MavenOptions.
func (in *MavenOptions) DeepCopy() *MavenOptions {
	if in == nil {
		return nil
	}
	out := new(MavenOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MavenRepository) DeepCopyInto(out *MavenRepository) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Layout != nil {
		in, out := &in.Layout, &out.Layout
		*out = new(string)
		**out = **in
	}
	if in.ReleasesEnabled != nil {
		in, out := &in.ReleasesEnabled, &out.ReleasesEnabled
		*out = new(bool)
		**out = **in
	}
	if in.ReleasesUpdatePolicy != nil {
		in, out := &in.ReleasesUpdatePolicy, &out.ReleasesUpdatePolicy
		*out = new(string)
		**out = **in
	}
	if in.SnapshotsEnabled != nil {
		in, out := &in.SnapshotsEnabled, &out.SnapshotsEnabled
		*out = new(bool)
		**out = **in
	}
	if in.SnapshotsUpdatePolicy != nil {
		in, out := &in.SnapshotsUpdatePolicy, &out.SnapshotsUpdatePolicy
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MavenRepository.
func (in *MavenRepository) DeepCopy() *MavenRepository {
	if in == nil {
		return nil
	}
	out := new(MavenRepository)
	in.DeepCopyInto(out)
	return out
}
