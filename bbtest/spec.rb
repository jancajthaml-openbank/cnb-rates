require 'turnip/rspec'
require 'json'
require 'thread'
require 'timeout'

Thread.abort_on_exception = true

RSpec.configure do |config|
  config.raise_error_for_unimplemented_steps = true
  config.color = true

  Dir.glob("./helpers/*_helper.rb") { |f| load f }
  config.include EventuallyHelper, :type => :feature
  Dir.glob("./steps/*_steps.rb") { |f| load f, true }

  config.before(:suite) do |_|
    print "[ suite starting ]\n"

    CNBHelper.start()

    ["/data", "/reports"].each { |folder|
      FileUtils.mkdir_p folder
      %x(rm -rf #{folder}/*)
    }

    $wall_instance_counter = 0
    $tenant_id = nil

    print "[ suite started  ]\n"
  end

  config.after(:suite) do |_|
    print "\n[ suite ending   ]\n"

    get_containers = lambda do |image|
      containers = %x(docker ps -a | awk '{ print $1,$2 }' | grep #{image} | awk '{print $1 }' 2>/dev/null)
      return ($? == 0 ? containers.split("\n") : [])
    end

    teardown_container = lambda do |container|
      label = %x(docker inspect --format='{{.Name}}' #{container})
      return unless $? == 0

      %x(docker exec #{container} systemctl stop cnb-rates.service 2>&1)
      %x(docker exec #{container} journalctl -o short-precise -u cnb-rates.service --no-pager >/reports/#{label.strip}.log 2>&1)
      %x(docker rm -f #{container} &>/dev/null || :)
    end

    capture_journal = lambda do |container|
      label = %x(docker inspect --format='{{.Name}}' #{container})
      return unless $? == 0

      %x(docker exec #{container} journalctl -o short-precise -u cnb-rates.service --no-pager >/reports/#{label.strip}.log 2>&1)
    end

    kill = lambda do |container|
      label = %x(docker inspect --format='{{.Name}}' #{container})
      return unless $? == 0
      %x(docker rm -f #{container.strip} &>/dev/null || :)
    end

    begin
      Timeout.timeout(5) do
        get_containers.call("openbank/cnb-rates").each { |container|
          teardown_container.call(container)
        }
      end
    rescue Timeout::Error => _
      get_containers.call("openbank/cnb-rates").each { |container|
        capture_journal.call(container)
        kill.call(container)
      }
      print "[ suite ending   ] (was not able to teardown container in time)\n"
    end

    CNBHelper.stop()

    print "[ suite cleaning ]\n"

    ["/data"].each { |folder|
      %x(rm -rf #{folder}/*)
    }

    print "[ suite ended    ]"
  end

end
