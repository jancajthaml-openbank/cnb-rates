require 'json'
require 'date'
require 'json-schema'
require 'thread'
require_relative '../shims/harden_webrick'

class CNBGetDailyMainRates < WEBrick::HTTPServlet::AbstractServlet

  def do_GET(request, response)
    status, body = process(request)

    response.status = status
    response.body = body
  end

  def process(request)
    query = request.query()

    if query.has_key?('date')
      for_date = Date.strptime(query['date'], "%d+%m+%Y")
    else
      for_date = Date.now
    end

    return 200, [
      (for_date.strftime('%d %b %Y') + " \#" + for_date.strftime('%j').to_i.to_s),
      "Country|Currency|Amount|Code|Rate",
      "Israel|shekel|1|ILS|10",
      "Japan|yen|100|JPY|10"
    ].join("\n")
  end
end

class CNBGetOtherFXRates < WEBrick::HTTPServlet::AbstractServlet

  def do_GET(request, response)
    status, body = process(request)

    response.status = status
    response.body = body
  end

  def process(request)
    query = request.query()

    if query.has_key?('month') && query.has_key?('year')
      for_date = Date.strptime(query['month'] + "." + query['year'], "%m.%Y").next_month.prev_day
    else
      for_date = Date.now
    end

    return 200, [
      (for_date.strftime('%d %b %Y') + " \#" + for_date.strftime('%j').to_i.to_s),
      "Country|Currency|Amount|Code|Rate",
      "Afghanistan|afghani|100|AFN|10",
      "Qatar|rial|1|QAR|10"
    ].join("\n")
  end
end

module CNBHelper

  def self.start
    self.server = nil

    begin
      self.server = WEBrick::HTTPServer.new(
        Port: 4000,
        Logger: WEBrick::Log.new("/dev/null"),
        AccessLog: [],
        SSLEnable: true
      )

    rescue Exception => err
      raise err
      raise "Failed to allocate server binding! #{err}"
    end

    self.server.mount "/en/financial_markets/foreign_exchange_market/exchange_rate_fixing/daily.txt", CNBGetDailyMainRates
    self.server.mount "/en/financial_markets/foreign_exchange_market/other_currencies_fx_rates/fx_rates.txt", CNBGetOtherFXRates

    self.server_daemon = Thread.new do
      self.server.start()
    end
  end

  def self.stop
    self.server.shutdown() unless self.server.nil?
    begin
      self.server_daemon.join() unless self.server_daemon.nil?
    rescue
    ensure
      self.server_daemon = nil
      self.server = nil
    end
  end

  class << self
    attr_accessor :server_daemon, :server
  end

end
