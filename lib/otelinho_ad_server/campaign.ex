defmodule OtelinhoAdServer.Campaign do
  use Ecto.Schema
  import Ecto.Changeset

  schema "campaigns" do
    field :budget, :decimal
    field :creative, :string
    field :start_date, :utc_datetime
    field :end_date, :utc_datetime
    field :goal, :integer
    field :max_bid, :decimal
    field :remaining_budget, :decimal
  end

  @doc false
  def changeset(campaign, attrs) do
    campaign
    |> cast(attrs, [:creative, :start_date, :end_date, :goal, :max_bid, :budget, :remaining_budget])
    |> validate_required([:creative, :start_date, :end_date, :goal, :max_bid, :budget, :remaining_budget])
  end
end
