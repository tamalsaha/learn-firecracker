sudo ip tuntap add tap2 mode tap
sudo ip link set tap2 up

sudo ip tuntap add tap3 mode tap
sudo ip addr add 172.26.0.2/24 dev tap3
sudo ip link set tap3 up

sudo sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward"
sudo iptables -t nat -A POSTROUTING -o bond0 -j MASQUERADE
sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
sudo iptables -A FORWARD -i tap3 -o bond0 -j ACCEPT


$ firectl \
  --kernel=./vmlinux \
  --root-drive=./root-drive-with-ssh.img \
  --kernel-opts="console=ttyS0 noapic reboot=k panic=1 pci=off nomodules rw" \
  --tap-device=tap2/AA:FC:ac:1a:00:02 \
  --tap-device=tap3/AA:FC:ac:1a:00:03

# MMDS
ip link set eth0 up
ip route add 169.254.169.254 dev eth0

# guest
ip addr add 172.26.0.3/24 dev eth1
ip link set eth1 up
ip route add default via 172.26.0.2 dev eth1


EgressInterface: bond0
instanceID: 0
ip0: 172.26.0.2 AA:FC:ac:1a:00:02
ip1: 172.26.0.3 AA:FC:ac:1a:00:03
tap2 tap3

# https://linuxlink.timesys.com/docs/static_ip
IP: 172.26.0.3::172.26.0.2:255.255.255.0::eth1:off:1.1.1.1:8.8.8.8:

