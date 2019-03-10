require 'date'

module TimeshiftHelper

  def self.get_dates_between(from, to)
    (from..to).select { |d| (1..5).include?(d.wday) }
  end

  def self.get_months_between(from, to)
    result = []

    df = Date.new(from.year, from.month, 1)
    dt = Date.new(to.year, to.month, 1)
    d = df
    while d < dt
      result << d.clone
      d = d.next_month
    end

    return result
  end

  def get_dates_between(*args)
    TimeshiftHelper.get_dates_between(*args)
  end

  def get_months_between(*args)
    TimeshiftHelper.get_months_between(*args)
  end

end
