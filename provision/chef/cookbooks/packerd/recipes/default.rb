
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
