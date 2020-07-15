// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackingStorageSpec) DeepCopyInto(out *BackingStorageSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackingStorageSpec.
func (in *BackingStorageSpec) DeepCopy() *BackingStorageSpec {
	if in == nil {
		return nil
	}
	out := new(BackingStorageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nfs) DeepCopyInto(out *Nfs) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nfs.
func (in *Nfs) DeepCopy() *Nfs {
	if in == nil {
		return nil
	}
	out := new(Nfs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Nfs) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NfsList) DeepCopyInto(out *NfsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Nfs, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NfsList.
func (in *NfsList) DeepCopy() *NfsList {
	if in == nil {
		return nil
	}
	out := new(NfsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NfsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NfsSpec) DeepCopyInto(out *NfsSpec) {
	*out = *in
	out.BackingStorage = in.BackingStorage
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NfsSpec.
func (in *NfsSpec) DeepCopy() *NfsSpec {
	if in == nil {
		return nil
	}
	out := new(NfsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NfsStatus) DeepCopyInto(out *NfsStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NfsStatus.
func (in *NfsStatus) DeepCopy() *NfsStatus {
	if in == nil {
		return nil
	}
	out := new(NfsStatus)
	in.DeepCopyInto(out)
	return out
}