require 'date'
require 'json-diff'

step "all CNB data are eventually synchronized" do ||
  date_from  = Date.parse('1991-01-01')
  date_to = if defined? @timeshift then @timeshift else Date.today end

  expected_days = TimeshiftHelper.get_dates_between(date_from, date_to).map { |d| d.strftime("%d.%m.%Y") }
  expected_months = TimeshiftHelper.get_months_between(date_from, date_to).map { |d| d.strftime("%m.%Y") }

  eventually(timeout: 30) {
    actual = Dir["/data/rates/cnb/raw/daily/fx-main/*"].select{ |f| File.file? f }.map{ |f| File.basename f }
    diff = JsonDiff.diff(actual, expected_days).select{ |item| item["op"] != "remove" and item["op"] != "move" }.map{ |item| item["value"] or item }
    raise "expectation failure:\ngot:\n#{actual}\nexpected:\n#{expected_days}\ndiff:\n#{diff}" if diff != []
  }

  eventually(timeout: 30) {
    actual = Dir["/data/rates/cnb/raw/monthly/fx-other/*"].select{ |f| File.file? f }.map{ |f| File.basename f }
    diff = JsonDiff.diff(actual, expected_months).select{ |item| item["op"] != "remove" and item["op"] != "move" }.map{ |item| item["value"] or item }
    raise "expectation failure:\ngot:\n#{actual}\nexpected:\n#{expected_months}\ndiff:\n#{diff}" if diff != []
  }
end
