defmodule OtelinhoAdServer.Index do
  import Ecto.Query, only: [from: 2]

  alias OtelinhoAdServer.Campaign
  alias OtelinhoAdServer.Repo

  def retrieve_active_campaigns() do
    now = DateTime.utc_now()
    query = from c in Campaign, where: c.start_date <= ^now and c.end_date >= ^now
    Repo.all(query)
  end
end
