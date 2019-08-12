

class BussinessLogic(object):

  def get_daily_main_fx(self, forDate):
    return '\n'.join([
      '{} #{}'.format(forDate.strftime('%d %b %Y'), int(forDate.strftime('%j'))),
      "Country|Currency|Amount|Code|Rate",
      "Israel|shekel|1|ILS|10",
      "Japan|yen|100|JPY|10"
    ])
