require 'date'
require 'json-diff'

step "all fx-main CNB data are downloaded" do ||
  date_from  = Date.parse('1991-01-01')
  date_to = if defined? @timeshift then @timeshift else Date.today end

  expected_days = TimeshiftHelper.get_dates_between(date_from, date_to).map { |d| d.strftime("%d.%m.%Y") }

  eventually(timeout: 60) {
    actual = Dir["/data/rates/cnb/raw/fx-main/*"].select{ |f| File.file? f }.map{ |f| File.basename f }
    diff = JsonDiff.diff(actual, expected_days).select { |item| item["op"] != "remove" and item["op"] != "move" }.map { |item| item["value"] or item }
    raise "expectation failure:\ngot:\n#{actual}\nexpected:\n#{expected_days}\ndiff:\n#{diff}" if diff != []
  }

  send ":operation unit :unit", "stop", "cnb-rates-import.service"
end

step "all fx-main CNB data are processed" do ||
  send "journalctl of :unit contains following", "cnb-rates-batch.service", ">>> Start <<<"
  send "journalctl of :unit contains following", "cnb-rates-batch.service", ">>> Stop <<<"

  date_from  = Date.parse('1991-01-01')
  date_to = if defined? @timeshift then @timeshift else Date.today end

  expected_days = TimeshiftHelper.get_dates_between(date_from, date_to).map { |d| d.strftime("%d.%m.%Y") }

  eventually(timeout: 60) {
    actual = Dir["/data/rates/cnb/processed/fx-main/d/*"].select{ |f| File.file? f }.map{ |f| File.basename f }
    diff = JsonDiff.diff(actual, expected_days).select { |item| item["op"] != "remove" and item["op"] != "move" }.map { |item| item["value"] or item }
    raise "expectation failure:\ngot:\n#{actual}\nexpected:\n#{expected_days}\ndiff:\n#{diff}" if diff != []
  }
end
