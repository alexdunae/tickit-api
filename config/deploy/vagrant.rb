server '127.0.0.1', {
  user: 'tickit',
  roles: %w{web app db},
  ssh_options: {
     keys: '/home/alex/.vagrant.d/insecure_private_key',
     port: 2222
  }
}
