# config valid only for Capistrano 3.1
lock '3.1.0'

set :application, 'tickit-api'
set :repo_url, 'git@github.com:alexdunae/tickit-api.git'
set :go_path, '/srv/www/tickit-api'

# Default branch is :master
# ask :branch, proc { `git rev-parse --abbrev-ref HEAD`.chomp }

# Default deploy_to directory is /var/www/my_app
set :deploy_to, '/srv/www/tickit-api'

# Default value for :log_level is :debug
# set :log_level, :debug

# Default value for :pty is false
# set :pty, true

# Default value for :linked_files is []
# set :linked_files, %w{config/database.yml}

# Default value for default_env is {}
#set :default_env, {
#  path: "/srv/www/tickit-api/bin:$PATH"
#  gopath: "/srv/www/tickit-api"
#}

namespace :deploy do
  desc 'Restart application'
  task :restart do
    on roles(:app), in: :sequence, wait: 5 do
      # Your restart mechanism here, for example:
      # execute :touch, release_path.join('tmp/restart.txt')
    end
  end

  after :publishing, :restart
end



namespace :go do
  task :build do
    on roles(:app), in: :sequence, wait: 5 do
    #within release_path do
      #with_env('GOPATH', release_path) do
        #execute "cd '#{release_path}'; #{fetch(:composer_command)} install"
        execute "GOPATH=#{release_path}; cd '#{release_path}/tickit-api'; go get"

        #run "go get" #
        #run "mkdir #{release_path}/bin"
        #run "cp /home/#{user}/go/bin/APPNAME #{release_path}/bin/"
      #end
    end
  end
end


after 'deploy:published', 'go:build'
