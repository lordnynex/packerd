deb = '/tmp/chefdk.deb'

remote_file deb do
  source 'https://omnitruck.chef.io/stable/chefdk/download?p=ubuntu&m=x86_64&pv=14.04&v=latest'
end

dpkg_package "chefdk" do
  source "/tmp/chefdk.deb"
  action :install
end

remote_file deb do
  action :delete
end

directory '/root/.berkshelf' do
  owner 'root'
  mode '0755'
  action :create
end

file "/root/.berkshelf/config.json" do
  owner 'root'
  mode '0755'
  content '{"ssl": { "verify": false }}'
  action :create
end

directory '/etc/supervisor' do
  owner 'root'
  group 'root'
  recursive true
end

template '/bin/docker_run.sh' do
  source 'docker_run.sh.erb'
  owner 'root'
  group 'root'
  mode '0755'
end

template '/etc/supervisor/supervisord.conf' do
  source 'supervisord.conf.erb'
  owner 'root'
  group 'root'
  mode '0755'
end
