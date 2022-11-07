
# instance = 0
sudo ip tuntap add tap1 mode tap
sudo ip link set tap1 up

sudo ip tuntap add tap2 mode tap
sudo ip addr add 172.26.0.1/30 dev tap2
sudo ip link set tap2 up

sudo sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward"
sudo iptables -t nat -A POSTROUTING -o bond0 -j MASQUERADE
sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
sudo iptables -A FORWARD -i tap2 -o bond0 -j ACCEPT

ssh -i id_m1 runner@172.26.0.2

# instance = 1
sudo ip tuntap add tap5 mode tap
sudo ip link set tap5 up

sudo ip tuntap add tap6 mode tap
sudo ip addr add 172.26.0.5/30 dev tap6
sudo ip link set tap6 up

sudo sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward"
sudo iptables -t nat -A POSTROUTING -o bond0 -j MASQUERADE
sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
sudo iptables -A FORWARD -i tap6 -o bond0 -j ACCEPT

# ssh -i id_m1 runner@172.26.0.6

ip0: 172.26.0.5 AA:FC:ac:1a:00:05
ip1: 172.26.0.6 AA:FC:ac:1a:00:06
tap5 tap6

