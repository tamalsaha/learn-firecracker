
# check mmds
curl -s --unix-socket /tmp/FCGoSDKSnapshotExample1022047977/fc-0.create http://localhost/mmds

# stop
curl --unix-socket /tmp/FCGoSDKSnapshotExample2494494834/fc-0.create -i \
  -X PUT 'http://localhost/actions'       \
  -H  'Accept: application/json'          \
  -H  'Content-Type: application/json'    \
  -d '{
      "action_type": "SendCtrlAltDel"
   }'


/tmp/FCGoSDKSnapshotExample1495375699/fc-0

curl --unix-socket /tmp/FCGoSDKSnapshotExample1998987509/fc-0.create -i \
    -X PUT "http://localhost/mmds"            \
    -H "Content-Type: application/json"       \
    -d '{
            "latest": {
                  "meta-data": {
                       "ami-id": "ami-12345678",
                       "reservation-id": "r-fea54097"
                  }
            }
    }'

curl -s -H "Accept: application/json" "http://169.254.169.254/latest"

ds=nocloud-net;s=http://169.254.169.254/latest/ network-config=dmVyc2lvbjogMgpldGhlcm5ldHM6CiAgZXRoMDoKICAgIGFkZHJlc3NlczoKICAgIC0gMTY5LjI1NC4xNjkuMjU0LzE2CiAgICBtYXRjaDoKICAgICAgbWFjYWRkcmVzczogQUE6RkM6YWM6MWE6MDA6MDIKICBldGgxOgogICAgYWRkcmVzc2VzOgogICAgLSAxNzIuMjYuMC4zLzI0CiAgICBnYXRld2F5NDogMTcyLjI2LjAuMgogICAgbWF0Y2g6CiAgICAgIG1hY2FkZHJlc3M6IEFBOkZDOmFjOjFhOjAwOjAzCiAgICBuYW1lc2VydmVyczoKICAgICAgYWRkcmVzc2VzOgogICAgICAtIDEuMS4xLjEKICAgICAgLSA4LjguOC44Cg==

https://askubuntu.com/questions/1104285/how-do-i-reload-network-configuration-with-cloud-init

https://askubuntu.com/questions/1405294/ubuntu-20-04-cloud-init-wont-configure-network


- https://netplan.io/troubleshooting


DI_LOG=stderr /usr/lib/cloud-init/ds-identify --force
rm /etc/cloud/cloud-init.disabled
yes | unminimize
cloud-init clean --logs
cloud-init init --local
# sudo cloud-init init
netplan generate
netplan apply

/etc/netplan/***

cloud-init static route
[route]
/run/systemd/network
https://manpages.ubuntu.com/manpages/bionic/man5/systemd.network.5.html

root@ubuntu-fc-uvm:~# ip r
default via 172.26.0.2 dev eth1
172.26.0.0/24 dev eth1 proto kernel scope link src 172.26.0.3
root@ubuntu-fc-uvm:~# netplan apply

root@ubuntu-fc-uvm:~# ip r
default via 172.26.0.2 dev eth1
default via 172.26.0.2 dev eth1 proto static
169.254.0.0/16 dev eth0 proto kernel scope link src 169.254.169.254
172.26.0.0/24 dev eth1 proto kernel scope link src 172.26.0.3



root@ubuntu-fc-uvm:~# ip r
default via 172.26.0.2 dev eth1
169.254.169.254 dev eth0 scope link
172.26.0.0/24 dev eth1 proto kernel scope link src 172.26.0.3

