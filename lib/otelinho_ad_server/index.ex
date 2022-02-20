defmodule OtelinhoAdServer.Index do
  use Agent

  import Ecto.Query, only: [from: 2]

  alias OtelinhoAdServer.Campaign
  alias OtelinhoAdServer.Repo

  def start_link(_opts) do
    Agent.start_link(fn -> [] end, name: __MODULE__)
  end

  def retrieve_active_campaigns() do
    Agent.get(__MODULE__, fn state -> state end)
  end

  def set_active_campaigns() do
    now = DateTime.utc_now()
    query = from c in Campaign, where: c.start_date <= ^now and c.end_date >= ^now
    campaigns = Repo.all(query)

    Agent.update(__MODULE__, fn _ -> campaigns end)
  end
end
