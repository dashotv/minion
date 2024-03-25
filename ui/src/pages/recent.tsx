import { useState } from "react";
import { JobsList } from "components/jobs";
import { Helmet } from "react-helmet-async";

const Recent = () => {
  // limit, skip, queries, etc
  const [page] = useState(1);
  const handleCancel = (id: string) => {
    console.log("cancel", id);
  };

  return (
    <>
      <Helmet>
        <title>Minion - Jobs</title>
        <meta name="description" content="runic" />
      </Helmet>
      <JobsList {...{ page, handleCancel }} />
    </>
  );
};

export default Recent;
