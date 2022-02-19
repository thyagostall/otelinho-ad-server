defmodule OtelinhoAdServer.Repo.Migrations.CreateCampaigns do
  use Ecto.Migration

  def change do
    create table(:campaigns) do
      add :creative, :string, null: false
      add :start_date, :utc_datetime, null: false
      add :end_date, :utc_datetime, null: false
      add :goal, :integer
      add :max_bid, :decimal
      add :budget, :decimal
      add :remaining_budget, :decimal
    end
  end
end
