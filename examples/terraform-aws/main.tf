variable "disk_image" {
  type = string
  default = "ami-0c1bea58988a989155"
}

variable "state" {
  type = string
  default = "running"
}

variable "machine_name" {
  type = string
  default = "devpod-devpod"
}

variable "disk_size" {
  type = string
  default = "50"
}

variable "instance_type" {
  type = string
  default = "t2.micro"
}

variable "vpc" {
  type = string
  default = ""
}

variable "region" {
  type = string
  default = ""
}

variable "ssh_key" {
  type = string
  default = "invalid"
}

provider "aws" {
  region = "${var.region}"
}


resource "aws_security_group" "devpod" {
  name = "devpod"
  description = "devpod"
  vpc_id = "${var.vpc}"

  // To Allow SSH Transport
  ingress {
    from_port = 22
    protocol = "tcp"
    to_port = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }
}


resource "aws_instance" "devpod" {
  ami = "${var.disk_image}"
  instance_type = "${var.instance_type}"

  vpc_security_group_ids = [
    aws_security_group.devpod.id
  ]
  root_block_device {
    delete_on_termination = true
    volume_size = "${var.disk_size}"
  }
  tags = {
    Name ="devpod"
    devpod = "${var.machine_name}"
  }

  user_data = <<EOF
#!/bin/sh
useradd devpod -d /home/devpod
mkdir -p /home/devpod
if grep -q sudo /etc/groups; then
  usermod -aG sudo devpod
elif grep -q wheel /etc/groups; then
  usermod -aG wheel devpod
fi
echo "devpod ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/91-devpod
mkdir -p /home/devpod/.ssh
echo "${var.ssh_key}" >> /home/devpod/.ssh/authorized_keys
chmod 0700 /home/devpod/.ssh
chmod 0600 /home/devpod/.ssh/authorized_keys
chown -R devpod:devpod /home/devpod
echo "${var.state}"
EOF

  depends_on = [ aws_security_group.devpod ]
}

output "public_ip" {
  value = aws_instance.devpod.public_ip
}
