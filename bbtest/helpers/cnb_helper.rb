require 'thread'
require 'agoo'

class YearlyHandler
  def call(req)
    qs = Hash.new
    req['QUERY_STRING'].split('&').each { |pair|
      key, value = pair.split('=', 2)
      qs[key] = value
    }

    if qs.has_key?('year')
      resp = [
        "Date|1 ILS|100 JPY"
      ]
    else
      resp = []
    end

    # fixme generate

    [ 200, { }, [ resp.join("\n") ] ]
  end
end

class MonthlyHandler
  def call(req)
    qs = Hash.new
    req['QUERY_STRING'].split('&').each { |pair|
      key, value = pair.split('=', 2)
      qs[key] = value
    }

    if qs.has_key?('month') && qs.has_key?('year')
      date = Date.strptime(qs['month'] + "." + qs['year'], "%m.%Y").next_month.prev_day
    else
      date = Date.now
    end

    resp = [
      (date.strftime('%d.%b %Y') + " \#" + date.strftime('%j').to_i.to_s),
      "Country|Currency|Amount|Code|Rate",
      "Afghanistan|afghani|100|AFN|10",
      "Qatar|rial|1|QAR|10"
    ]

    [ 200, { }, [ resp.join("\n") ] ]
  end
end

class DailyHandler
  def call(req)
    qs = Hash.new
    req['QUERY_STRING'].split('&').each { |pair|
      key, value = pair.split('=', 2)
      qs[key] = value
    }

    if qs.has_key?('date')
      date = Date.strptime(qs['date'], "%d.%m.%Y")
    else
      date = Date.now
    end

    resp = [
      (date.strftime('%d.%b %Y') + " \#" + date.strftime('%j').to_i.to_s),
      "Country|Currency|Amount|Code|Rate",
      "Israel|shekel|1|ILS|10",
      "Japan|yen|100|JPY|10"
    ]

    [ 200, { }, [ resp.join("\n") ] ]
  end
end

module CNBHelper

  def self.start
    begin
      Agoo::Server.init(8080, 'root')
    rescue Exception => _
      raise "Failed to allocate server binding!"
    end

    Agoo::Server.handle(:GET, "/en/financial_markets/foreign_exchange_market/exchange_rate_fixing/daily.txt", DailyHandler.new)
    Agoo::Server.handle(:GET, "/en/financial_markets/foreign_exchange_market/other_currencies_fx_rates/fx_rates.txt", MonthlyHandler.new)
    Agoo::Server.handle(:GET, "/en/financial_markets/foreign_exchange_market/exchange_rate_fixing/year.txt", YearlyHandler.new)

    self.server_daemon = Thread.new do
      Agoo::Server.start()
    end
  end

  def self.stop
    Agoo::Server.shutdown()
    begin
      self.server_daemon.join() unless self.server_daemon.nil?
    rescue
    ensure
      self.server_daemon = nil
    end
  end

  class << self
    attr_accessor :server_daemon
  end

end
