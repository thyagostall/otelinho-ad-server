defmodule OtelinhoAdServer.Repo.Migrations.CreateBeacons do
  use Ecto.Migration

  def change do
    create table(:beacons) do
      add :campaign_id, :integer
      add :event, :string
      add :timestamp, :time
    end
  end
end
