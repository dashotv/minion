import { Button, Stack } from "@mui/material";
import { Pill } from "components/common";
import { Stats } from "./types";
import { useNavigate } from "react-router-dom";

export const JobsStats = ({ stats }: { stats: Stats }) => {
  const navigate = useNavigate();

  if (!stats) {
    return null;
  }
  return (
    <Stack direction="row" spacing={0} justifyContent="end">
      <Button onClick={() => navigate("search?status=pending")}>
        <Pill name="Pending" color="gray" value={stats.pending || 0} />
      </Button>
      <Button onClick={() => navigate("search?status=queued")}>
        <Pill name="Queued" color="secondary" value={stats.queued || 0} />
      </Button>
      <Button onClick={() => navigate("search?status=running")}>
        <Pill name="Running" color="primary" value={stats.running || 0} />
      </Button>
      <Button onClick={() => navigate("search?status=cancelled")}>
        <Pill name="Cancelled" color="warning" value={stats.cancelled} />
      </Button>
      <Button onClick={() => navigate("search?status=failed")}>
        <Pill name="Failed" color="error" value={stats.failed} />
      </Button>
    </Stack>
  );
};
