import { getPublicationById } from "@/services/publications";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";

const Publication = () => {
  const { id } = useParams();
  const { data, isLoading, error } = useQuery({
    queryKey: ["publications"],
    queryFn: () => getPublicationById(id!),
  });

  if (isLoading) return "Loading...";
  if (error) return "An error has occurred: " + error.message;

  return (
    <div>
      <h1>{data?.title}</h1>
      <p>{data?.description}</p>
      <p>{data?.rating}</p>
      <p>{data?.image}</p>
      <p>{data?.created_at}</p>
      <p>{data?.updated_at}</p>
    </div>
  );
};

export default Publication;
