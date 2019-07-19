require_relative 'placeholders'

step "cnb-rates is restarted" do ||
  ids = %x(systemctl -t service --no-legend | awk '{ print $1 }')
  expect($?).to be_success, ids

  ids = ids.split("\n").map(&:strip).reject { |x|
    x.empty? || !x.start_with?("cnb-rates-")
  }.map { |x| x.chomp(".service") }

  expect(ids).not_to be_empty

  ids.each { |e|
    %x(systemctl restart #{e} 2>&1)
  }

  eventually() {
    ids.each { |e|
      out = %x(systemctl show -p SubState #{e} 2>&1 | sed 's/SubState=//g')
      expect(out.strip).to eq("running")
    }
  }
end

step "cnb-rates is running with mocked CNB Gateway" do ||
  params = [
    "CNB_GATEWAY=https://127.0.0.1:4000",
  ].join("\n")

  ts = if defined? @timeshift then @timeshift else Date.today end
  formatted = ts.strftime("%Y-%m-%d %H:%M:%S")

  %x(timedatectl set-ntp 0)
  expect($?).to be_success, "failed to set-ntp"
  %x(timedatectl set-local-rtc 0)
  expect($?).to be_success, "failed to set-local-rtc"
  %x(timedatectl set-time "#{formatted}")
  expect($?).to be_success, "failed to set time to #{formatted}"
  %x(systemctl restart cron)
  send "cnb-rates is reconfigured with", params
end

step "cnb-rates is reconfigured with" do |configuration|
  params = Hash[configuration.split("\n").map(&:strip).reject(&:empty?).map { |el| el.split '=' }]
  config = Array[UnitHelper.default_config.merge(params).map {|k,v| "CNB_RATES_#{k}=#{v}"}]
  config = config.join("\n").inspect.delete('\"')

  %x(mkdir -p /etc/init)
  %x(echo '#{config}' > /etc/init/cnb-rates.conf)

  ids = %x(systemctl list-units | awk '{ print $1 }')
  expect($?).to be_success, ids

  ids = ids.split("\n").map(&:strip).reject { |x|
    x.empty? || !x.start_with?("cnb-rates-")
  }.map { |x| x.chomp(".service") }

  ids.each { |e|
    %x(systemctl restart #{e} 2>&1)
  }

  eventually() {
    ids.each { |e|
      out = %x(systemctl show -p SubState #{e} 2>&1 | sed 's/SubState=//g')
      expect(out.strip).to eq("running"), "#{e} is #{out.strip} expected running"
    }
  }
end
