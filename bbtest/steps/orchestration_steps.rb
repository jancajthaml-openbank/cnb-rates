require 'tempfile'

step "no :container :label is running" do |container, label|
  containers = %x(docker ps -a -f name=#{label} | awk '$2 ~ "#{container}" {print $1}' 2>/dev/null)
  expect($?).to be_success

  ids = containers.split("\n").map(&:strip).reject(&:empty?)
  return if ids.empty?

  ids.each { |id|
    eventually(timeout: 2) {
      puts "wanting to kill #{id}"
      send ":container running state is :state", id, false

      label = %x(docker inspect --format='{{.Name}}' #{id})
      label = ($? == 0 ? label.strip : id)

      %x(docker exec #{container} journalctl -o short-precise -u cnb-rates.service --no-pager >/reports/#{label}.log 2>&1)
      %x(docker rm -f #{id} &>/dev/null || :)
    }
  }
end

step ":container running state is :state" do |container, state|
  eventually(timeout: 3) {
    %x(docker exec #{container} systemctl stop cnb-rates.service 2>&1) unless state
    expect($?).to be_success

    %x(docker #{state ? "start" : "stop"} #{container} >/dev/null 2>&1)

    container_state = %x(docker inspect -f {{.State.Running}} #{container} 2>/dev/null)
    expect($?).to be_success
    expect(container_state.strip).to eq(state ? "true" : "false")
  }
end

step ":container :version is started with" do |container, version, label, params|
  containers = %x(docker ps -a --filter name=#{label} --filter status=running --format "{{.ID}} {{.Image}}")
  expect($?).to be_success
  containers = containers.split("\n").map(&:strip).reject(&:empty?)

  unless containers.empty?
    id, image = containers[0].split(" ")
    return if (image == "#{container}:#{version}")
  end

  send "no :container :label is running", container, label

  prefix = ENV.fetch('COMPOSE_PROJECT_NAME', "")
  my_id = %x(cat /etc/hostname).strip
  args = [
    "docker",
    "run",
    "-d",
    "--net=#{prefix}_default",
    "--volumes-from=#{my_id}",
    "--log-driver=json-file",
    "-h #{label}",
    "--net-alias=#{label}",
    "--name=#{label}",
    "--privileged"
  ] << params << [
    "#{container}:#{version}",
    "2>&1"
  ]

  id = %x(#{args.join(" ")})
  expect($?).to be_success, id

  eventually(timeout: 3) {
    send ":container running state is :state", id, true
  }
end

step "cnb-rates is running" do ||
  with_deadline(timeout: 5) {
    send ":container :version is started with", "openbank/cnb-rates", ENV.fetch("VERSION", "latest"), "cnb-rates", [
      "-v /sys/fs/cgroup:/sys/fs/cgroup:ro"
    ]
  }
end

step "cnb-rates is running with following configuration" do |configuration|
  with_deadline(timeout: 5) {
    send ":container :version is started with", "openbank/cnb-rates", ENV.fetch("VERSION", "latest"), "cnb-rates", [
      "-v /sys/fs/cgroup:/sys/fs/cgroup:ro"
    ]
  }

  params = configuration.split("\n").map(&:strip).reject(&:empty?).join("\n").inspect.delete('\"')

  containers = %x(docker ps -a --filter name=cnb-rates --filter status=running --format "{{.ID}} {{.Image}}")
  expect($?).to be_success
  containers = containers.split("\n").map(&:strip).reject(&:empty?)

  expect(containers).not_to be_empty

  id = containers[0].split(" ")[0]

  %x(docker exec #{id} bash -c "echo -e '#{params}' > /etc/init/cnb-rates.conf" 2>&1)
  %x(docker exec #{id} systemctl restart cnb-rates.service 2>&1)
end
