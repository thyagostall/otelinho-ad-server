defmodule OtelinhoAdServer.Pacing do
  use Agent

  def start_link(_opts) do
    Agent.start_link(fn -> %{} end, name: __MODULE__)
  end

  def get(campaign_id) do
    Agent.get(__MODULE__, fn state -> Map.get(state, campaign_id) end)
  end

  def put(campaign_id, value) do
    Agent.update(__MODULE__, fn state -> Map.put(state, campaign_id, value) end)
  end

  def filter(campaigns) do
    pacing_factors = Agent.get(__MODULE__, fn state -> state end)
    Enum.filter(campaigns, fn campaign -> Map.get(pacing_factors, campaign.id) < :rand.uniform(Integer.pow(2, 32)) end)
  end

  def populate(campaigns) do
    campaigns
    |> Enum.each(fn campaign -> put(campaign.id, Integer.pow(2, 31)) end)
  end
end
