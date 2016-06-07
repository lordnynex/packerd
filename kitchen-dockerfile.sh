#!/bin/bash

out=kitchen-dockerfile

echo "FROM dockyardaws.cloud.capitalone.com/tompscanlan/packerd" > $out
if [ -n $http_proxy ]; then echo "ENV http_proxy $http_proxy" >> $out; fi
if [ -n $https_proxy ]; then echo "ENV https_proxy $https_proxy" >> $out; fi
echo "ENV DEBIAN_FRONTEND noninteractive" >> $out

cat << EOF >> $out
RUN dpkg-divert --local --rename --add /sbin/initctl && \\
 ln -sf /bin/true /sbin/initctl && \\
 apt-get update && \\
 apt-get install -y sudo openssh-server curl lsb-release && \\
 useradd -d /home/kitchen -m -s /bin/bash kitchen && \\
 echo kitchen:kitchen | chpasswd && \\
 echo 'kitchen ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers && \\
 mkdir -p /etc/sudoers.d && \\
 echo 'kitchen ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers.d/kitchen && \\
 chmod 0440 /etc/sudoers.d/kitchen && \\
 mkdir -p /home/kitchen/.ssh && \\
EOF

 if [ -f .kitchen/docker_id_rsa.pub ]; then
   ssh_key=`cat .kitchen/docker_id_rsa.pub`
   echo -n "echo '$ssh_key' >> /home/kitchen/.ssh/authorized_keys ;" >> $out
fi

