sudo apt-get update
sudo apt-get -y install python-software-properties
sudo add-apt-repository -y ppa:duh/golang
sudo apt-get update
sudo apt-get -y install golang git mercurial mysql-server

echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.bashrc
source /home/vagrant/.bashrc
mkdir -p $GOPATH/src/tickit

cd $GOPATH/src/tickit
rm -rfv $GOPATH/src/tickit
git clone git@github.com:alexdunae/tickit-api.git $GOPATH/src/tickit --depth 1

sudo cp $GOPATH/src/tickit/tickit-api/checkin-api.conf /etc/tickit-api.conf
sudo cp $GOPATH/src/tickit/upstart.conf /etc/init/tickit-api.conf

cd $GOPATH/src/tickit/tickit-api
go get
