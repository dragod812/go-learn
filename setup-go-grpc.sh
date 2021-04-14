# install go
wget https://golang.org/dl/go1.15.linux-amd64.tar.gz
sudo tar xzvf go1.15.linux-amd64.tar.gz
sudo mv go /usr/local
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# make GOPATH and src
mkdir go && mkdir go/src

# install protobuf compiler and grpc
sudo apt install -y protobuf-compiler
go get -u google.golang.org/grpc

# Change terminal color
export PS1="\e[0;36m\W\$\e[m "
