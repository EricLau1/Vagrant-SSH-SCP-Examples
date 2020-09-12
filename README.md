# Exemplos com Vagrant

## Criando uma Máquina Virtual

1. Criar o arquivo `Vagrantfile`

```rb
Vagrant.configure("2") do |config|
	config.vm.box = "bento/ubuntu-18.04"

    config.vm.define :bucket do |bucket_config|

        bucket_config.vm.network "private_network", ip: "192.168.50.11"
        bucket_config.vm.hostname = 'bucket'

        bucket_config.vm.provider :virtualbox do |vb|
            vb.name = "bucket"
        end

    end

end
```

## Comandos Gerais

```bash
    # Validar o arquivo
    vagrant validate

    # Executa o Vagrantfile
    vagrant up

    # Destruindo uma vm
    vagrant destroy bucket

    # Verificando status
    vagrant status

    # Verificando status global
    vagrant global-status

    # Listando vms do virtual box
    vagrant box list

    # Entrar na vm sem precisar digitar a senha
    vagrant ssh

    # Entrar na vm com senha.
    # A senha padrão é: vagrant
    ssh vagrant@192.168.50.11

    # Copiando arquivos para a VM
    echo "Teste $(date)" >> test.txt

    scp test.txt vagrant@192.168.50.11:
```

## Acessar Máquina Virtual com Chave Pública

> Nesta parte já considera que a chave publica e privada tenha sido criada

1. Instalar o `Ansible`:

```bash
    sudo apt-get install ansible
```

2. Criar uma pasta `playbooks` e criar o arquivo `ssh-connections.yaml`:

```yml
- hosts: all
  tasks:
    - name: Criando Permissão com Chave Pública
      become: yes
      blockinfile:
        marker: ""
        marker_begin: ""
        marker_end: ""
        block: "{{ lookup('file', '../../ssh-keys/exemplo.pub') }}"
        dest: /home/vagrant/.ssh/authorized_keys
```

3. Editar o `Vagrantfile` para ler o arquivo `ssh-connections.yaml`:

```rb
Vagrant.configure("2") do |config|
	config.vm.box = "bento/ubuntu-18.04"

    config.vm.define :bucket do |bucket_config|

        bucket_config.vm.network "private_network", ip: "192.168.50.11"
        bucket_config.vm.hostname = 'bucket'

        bucket_config.vm.provider :virtualbox do |vb|
            vb.name = "bucket"
        end

        bucket_config.vm.provision "ansible" do |ansible|
            ansible.playbook = "playbooks/ssh-connections.yaml"
		end

    end
    
end
```

4. Executar os comandos:

```bash
    # Suspendendo a VM
    vagrant suspend bucket

    vagrant up

    vagrant provision    
```

5. Acessando a vm com a chave privada:

```bash
    ssh -i ssh-keys/exemplo vagrant@192.168.50.11

    # Dependendo do usuário que esteja tentando se conectar
    # será necessário executar o seguinte comando:
    ssh -o StrictHostKeyChecking=no -i ssh-keys/exemplo vagrant@192.168.50.11

    # Enviando varios arquivos para a VM
    # sem precisar de password
    scp -o StrictHostKeyChecking=no -i ssh-keys/exemplo  prog/files/*.pdf vagrant@192.168.50.11:

    # Copiando da VM para o diretorio Local
    scp -o StrictHostKeyChecking=no -i ssh-keys/exemplo vagrant@192.168.50.11:/home/vagrant/*.pdf .
```