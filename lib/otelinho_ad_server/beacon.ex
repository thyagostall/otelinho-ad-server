defmodule OtelinhoAdServer.Beacon do
  use Ecto.Schema
  import Ecto.Changeset

  schema "beacons" do
    field :campaign_id, :integer
    field :event, :string
    field :timestamp, :time
  end

  @doc false
  def changeset(beacon, attrs) do
    beacon
    |> cast(attrs, [:campaign_id, :event, :timestamp])
    |> validate_required([:campaign_id, :event, :timestamp])
  end
end
