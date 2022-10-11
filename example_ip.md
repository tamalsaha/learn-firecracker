```
# ip r
default via 10.168.0.1 dev eth0
10.168.0.0/24 dev eth0 proto kernel scope link src 10.168.0.15


# ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether ca:78:bc:64:c0:4b brd ff:ff:ff:ff:ff:ff
    inet 10.168.0.15/24 brd 10.168.0.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::c878:bcff:fe64:c04b/64 scope link
       valid_lft forever preferred_lft forever


# cat /proc/cmdline
ip=10.168.0.15::10.168.0.1:255.255.255.0::eth0:off:1.1.1.1:8.8.8.8: root=/dev/vda rw virtio_mmio.device=4K@0xd0000000:5 virtio_mmio.device=4K@0xd0001000:6
```

------------------------------

```
# ip r
default via 45.33.117.1 dev eth0 proto static
45.33.117.0/24 dev eth0 proto kernel scope link src 45.33.117.107
172.17.0.0/16 dev docker0 proto kernel scope link src 172.17.0.1 linkdown


# ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    link/ether f2:3c:93:77:7b:4f brd ff:ff:ff:ff:ff:ff
    inet 45.33.117.107/24 brd 45.33.117.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 2600:3c00::f03c:93ff:fe77:7b4f/64 scope global dynamic mngtmpaddr noprefixroute
       valid_lft 89sec preferred_lft 29sec
    inet6 fe80::f03c:93ff:fe77:7b4f/64 scope link
       valid_lft forever preferred_lft forever
3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default
    link/ether 02:42:00:c2:4e:2f brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
       valid_lft forever preferred_lft forever
    inet6 fe80::42:ff:fec2:4e2f/64 scope link
       valid_lft forever preferred_lft forever


# cat /proc/cmdline
BOOT_IMAGE=/boot/vmlinuz-5.15.0-48-generic root=/dev/sda ro console=ttyS0,19200n8 net.ifnames=0 scsi_mod.scan=sync
```
