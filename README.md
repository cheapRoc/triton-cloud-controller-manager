# Triton Cloud Controller Manager

_triton-cloud-controller-manager_ is an external [Kubernetes][kube] cloud
controller manager for automating Kubernetes nodes running on Joyent's Triton
cloud platform.

# Background

External cloud providers were introduced as an alpha feature in Kubernetes
release *1.6*. This repository contains an initial implementation of an external
cloud provider for running Kubernetes resources on top of Joyent's Triton
public/private cloud.

An external cloud provider acts like any other Kubernetes controller except is
responsible for cloud provider-specific control loops required for the
functioning of Kubernetes itself. These loops were originally part of the
`kube-controller-manager` daemon, but were tightly coupling the
`kube-controller-manager` to cloud-provider specific code. In order to free the
core Kubernetes project of this dependency, the concept of a
`cloud-controller-manager` was introduced.

`cloud-controller-manager` allows cloud vendors and Kubernetes core to evolve
independent of each other. In prior releases, the core Kubernetes code was
dependent upon cloud provider-specific code for functionality. In future
releases, code specific to cloud vendors should be maintained by the cloud
vendors themselves, and linked to `cloud-controller-manager` while running
Kubernetes.

In order to use this controller you must disable these internal controller loops
within the `kube-controller-manager` if you are running the
`triton-cloud-controller-manager`. You can disable the controller loops by
setting the `--cloud-provider` flag to `external` when starting the
`kube-controller-manager`.

# Features

The following controllers will be implemented by the `triton-cloud-controller-manager`:

* **WiP** A tool that allows storing `triton-go` authentication credentials
  inside Kubernetes as an internally held secret. This allows
  `triton-cloud-controller-manager` to access the Triton's CloudAPI
  at`triton-go` client integration points. Ref [make_secret.go][secret].

* **WIP** Node Controller: For checking Triton's CloudAPI to determine if a KVM
  node has been deleted in the cloud after it stops responding.

* **WIP** Route Controller: For setting up routes in the underlying cloud
  infrastructure, possibly subnets/networking.

* **WIP** Service Controller: For creating, updating and deleting cloud provider
  DNS load balancing by switching out CNS tags on given nodes.

* **WIP** Metrics Controller: For enabling and directing metrics traffics into
  Triton's CMON service.

# Developing

```sh
$ make help
```

# Running

**NOTE**: In order to test this daemon you'll probably want a Kubernetes cluster
  of your own. I followed Joyent's [Kubernetes The Hard Way][hard-way] blog
  article but I'm sure you can use whatever you've got running. Please let me
  know what doesn't work.

[hard-way]: https://www.joyent.com/blog/kubernetes-the-hard-way

First, create a separate SSH key (using `ssh-keygen`) and apply it to your
Triton account. This will be used to authenticate all `triton-go` requests
coming out of this controller.

Next, disable any existing cloud provider control loops running within
`kube-controller-manager` by setting the `--cloud-provider` flag to `external`
when starting the `kube-controller-manager` itself. If you used
[the hard way][hard-way] then you'll need to edit
`/etc/systemd/system/kube-controller-manager.service` to include the service
with the flag mentioned above.

You'll need to build the binary specifically for the platform you've deployed
Kubernetes onto. My default is (will be) `linux/amd64`.

```sh
$ make build
```

Once the binary is built, copy it to your deploy host and run the
following... or do we run it from `systemd`???

```sh
$ kubectl ...something or other....
```

More notes here...

- https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/
- https://github.com/coreos/coreos-kubernetes/blob/adcbe533/multi-node/aws/artifacts/manifests/controller/kube-controller-manager.yaml#L16

[kube]: https://kubernetes.io
[secret]: https://github.com/kubernetes/contrib/blob/master/ingress/controllers/gce/examples/https/make_secret.go
