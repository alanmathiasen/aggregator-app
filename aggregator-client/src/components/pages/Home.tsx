import { useQuery } from "@tanstack/react-query";
import { getAllPublications } from "@/services/publications";
import "../../App.css";

const Home = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["publications"],
    queryFn: getAllPublications,
  });

  if (isLoading) return "Loading...";
  if (error) return "An error has occurred: " + error.message;
  return (
    <div>
      {data?.publications.map((publication) => (
        <p key={publication.id}>
          {publication.title} {publication.id}
        </p>
      ))}
    </div>
  );
};

export default Home;
