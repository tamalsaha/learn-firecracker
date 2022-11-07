

modprobe: FATAL: Module ip_tables not found in directory /lib/modules/5.4.0-131-generic
iptables v1.8.4 (legacy): can't initialize iptables table `nat': Table does not exist (do you need to insmod?)
Perhaps iptables or your kernel needs to be upgraded.


curl -L https://gist.githubusercontent.com/tamalsaha/af2f99c80f84410253bd1e532bdfabc7/raw/1a5ea8aba0f06d71d0cf1934b2b92aa2b1b10528/build_rootfs.sh | bash -s -- bionic 20G




apt install --reinstall linux-modules-`uname -r`

ls -l /lib/modules/$(uname -r)/kernel/net/ipv4/netfilter/ip_tables.ko
