# -*- mode: ruby -*-

# Options {{{
#
APP_VM_NAME = "ldl"
APP_MEMORY = "1024"
APP_CPUS = "1"
#
# }}}

# Vagrant 2.0.x {{{
#
Vagrant.configure("2") do |config|

    config.vm.box = "ubuntu/trusty64"
    config.vm.box_check_update = false

    config.vm.provision :shell, :path => ".vagrant_bootstrap.sh"
    config.ssh.shell = "bash -c 'BASH_ENV=/etc/profile exec bash'"

    if Vagrant.has_plugin?("vagrant-cachier")
        config.cache.scope = :box

        config.cache.synced_folder_opts = {
            type: :nfs,
            mount_options: ['rw', 'vers=3', 'tcp', 'nolock']
        }
    end

    config.vm.define "master", primary: true do |master|
        master.vm.hostname = "ldl-master"
        master.vm.network :private_network, ip: "48.44.44.44"

        master.vm.provider :virtualbox do |vb|
            vb.name = "ldl-master"
            vb.customize ["modifyvm", :id, "--memory", APP_MEMORY]
            vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
            vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
            vb.customize ["modifyvm", :id, "--ioapic", "on"]
            vb.customize ["modifyvm", :id, "--cpus", APP_CPUS]

            # vb.customize ["createhd",  "--filename", "ldl_master_zfs_disk0", "--size", "2048"]
            # vb.customize ["storageattach", :id, "--storagectl", "SATAController", "--port", "1", "--type", "hdd", "--medium", "ldl_master_zfs_disk0.vdi"]
        end

    end

    config.vm.define "slave" do |slave|
        slave.vm.hostname = "ldl-slave"
        slave.vm.network :private_network, ip: "49.44.44.44"

        slave.vm.provider :virtualbox do |vb|
            vb.name = "ldl-slave"
            vb.customize ["modifyvm", :id, "--memory", APP_MEMORY]
            vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
            vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
            vb.customize ["modifyvm", :id, "--ioapic", "on"]
            vb.customize ["modifyvm", :id, "--cpus", APP_CPUS]

            # vb.customize ["createhd",  "--filename", "ldl_slave_zfs_disk0", "--size", "2048"]
            # vb.customize ["storageattach", :id, "--storagectl", "SATAController", "--port", "1", "--type", "hdd", "--medium", "ldl_slave_zfs_disk0.vdi"]
        end
    end
end
#
# }}}
