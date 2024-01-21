import { Publication } from "@/types/Publication";
import axios from "axios";

const URL = "http://localhost:8080/api/v1";

export const getAllPublications = async (): Promise<{
  publications: Publication[];
}> => {
  const response = await axios.get(`${URL}/publications`);
  return response.data;
};

export const getPublicationById = async (id: string): Promise<Publication> => {
  const response = await axios.get(`${URL}/publications/${id}`);
  return response.data;
};
