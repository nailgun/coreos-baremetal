
# Getting Started with Docker

In this tutorial, we'll run `bootcfg` on your Linux machine with Docker to network boot and provision a cluster of QEMU/KVM CoreOS machines locally. You'll be able to create Kubernetes clustes, etcd clusters, and test network setups.

*Note*: To provision physical machines, see [network setup](network-setup.md) and [deployment](deployment.md).

## Requirements

Install the package dependencies and start the Docker daemon.

    # Fedora
    sudo dnf install docker virt-install virt-manager
    sudo systemctl start docker

    # Debian/Ubuntu
    # check Docker's docs to install Docker 1.8+ on Debian/Ubuntu
    sudo apt-get install virt-manager virtinst qemu-kvm

Clone the [coreos-baremetal](https://github.com/coreos/coreos-baremetal) source which contains the examples and scripts.

    git clone https://github.com/coreos/coreos-baremetal.git
    cd coreos-baremetal

Download CoreOS image assets referenced by the `etcd-docker` [example](../examples) to `examples/assets`.

    ./scripts/get-coreos stable 1185.3.0 ./examples/assets

For development convenience, add `/etc/hosts` entries for nodes so they may be referenced by name as you would in production.

    # /etc/hosts
    ...
    172.17.0.21 node1.example.com
    172.17.0.22 node2.example.com
    172.17.0.23 node3.example.com

## Containers

Run the latest `bootcfg` Docker image from `quay.io/coreos/bootcfg` with the `etcd-docker` example. The container should receive the IP address 172.17.0.2 on the `docker0` bridge.

    sudo docker pull quay.io/coreos/bootcfg:latest
    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd:/var/lib/bootcfg/groups:Z quay.io/coreos/bootcfg:latest -address=0.0.0.0:8080 -log-level=debug

or run the latest tagged release.

    sudo docker run -p 8080:8080 --rm -v $PWD/examples:/var/lib/bootcfg:Z -v $PWD/examples/groups/etcd:/var/lib/bootcfg/groups:Z quay.io/coreos/bootcfg:v0.4.1 -address=0.0.0.0:8080 -log-level=debug

Take a look at the [etcd groups](../examples/groups/etcd) to get an idea of how machines are mapped to Profiles. Explore some endpoints exposed by the service, say for QEMU/KVM node1.

* iPXE [http://127.0.0.1:8080/ipxe?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/ipxe?mac=52:54:00:a1:9c:ae)
* Ignition [http://127.0.0.1:8080/ignition?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/ignition?mac=52:54:00:a1:9c:ae)
* Metadata [http://127.0.0.1:8080/metadata?mac=52:54:00:a1:9c:ae](http://127.0.0.1:8080/metadata?mac=52:54:00:a1:9c:ae)

## Network

Since the virtual network has no network boot services, use the `dnsmasq` image to create an iPXE network boot environment which runs DHCP, DNS, and TFTP.

    sudo docker run --name dnsmasq --cap-add=NET_ADMIN -v $PWD/contrib/dnsmasq/docker0.conf:/etc/dnsmasq.conf:Z quay.io/coreos/dnsmasq -d

In this case, dnsmasq runs a DHCP server allocating IPs to VMs between 172.17.0.43 and 172.17.0.99, resolves `bootcfg.foo` to 172.17.0.2 (the IP where `bootcfg` runs), and points iPXE clients to `http://bootcfg.foo:8080/boot.ipxe`.

## Client VMs

Create QEMU/KVM VMs which have known hardware attributes. The nodes will be attached to the `docker0` bridge, where Docker's containers run.

    sudo ./scripts/libvirt create-docker
    sudo virt-manager

You can use `virt-manager` to watch the console and reboot VM machines with

    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt start

## Verify

The VMs should network boot and provision themselves into a three node etcd cluster, with other nodes behaving as etcd proxies.

The example profile added autologin so you can verify that etcd works between nodes.

    systemctl status etcd2
    etcdctl set /message hello
    etcdctl get /message

## Cleanup

Clean up the containers and VM machines.

    sudo docker rm -f dnsmasq
    sudo ./scripts/libvirt poweroff
    sudo ./scripts/libvirt destroy

## Going Further

Learn more about [bootcfg](bootcfg.md) or explore the other [example](../examples) clusters. Try the [k8s example](kubernetes.md) to produce a TLS-authenticated Kubernetes cluster you can access locally with `kubectl`.

