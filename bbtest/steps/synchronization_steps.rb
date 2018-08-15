require 'date'

step "all CNB data are eventually synchronized" do ||
  eventually(timeout: 1) {
    expect(File.directory?("/data/cnb-rates/raw/fx-main/daily")).to be(true), "directory /data/cnb-rates/raw/fx-main/daily does not exists"
  }

  date_from  = Date.parse('1991-01-01')
  date_to = if defined? @timeshift then @timeshift else Date.today end

  last_date = date_to.next_month.prev_day.strftime("%d.%m.%Y")

  expected = (date_from..date_to).select { |d| (1..5).include?(d.wday) }.map { |d| d.strftime("%d.%m.%Y") }

  eventually(timeout: 60) {
    actual = Dir["/data/cnb-rates/raw/fx-main/daily/*"].select{ |f| File.file? f }.map{ |f| File.basename f }
    expect(actual).to satisfy { |v| v.length >= expected.length }, "#expected to found #{expected.length} files but found #{actual.length}"
    expect(actual).to include(*expected)
  }
end
