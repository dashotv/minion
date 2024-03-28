import axios from 'axios';

import { useQuery } from '@tanstack/react-query';

import { JobsResponse } from './types';

export const getJobsFor = async (id: string, page: number) => {
  const response = await axios.get(`/api/minion/jobs/?page=${page}&client=${id}`);
  return response.data.jobs as JobsResponse;
};

export const getJobs = async (page: number, status: string) => {
  const response = await axios.get(`/api/minion/jobs?limit=250&page=${page}&status=${status}`);
  return response.data as JobsResponse;
};

export const queueJob = async (name: string, client: string) => {
  const response = await axios.post(`/api/minion/jobs?job=${name}&client=${client}`);
  return response.data;
};

export const deleteJob = async (id: string, hard: boolean) => {
  const response = await axios.delete(`/api/minion/jobs/${id}?hard=${hard}`);
  return response.data;
};

export const useJobsQuery = (page: number, status: string) =>
  useQuery({
    queryKey: ['jobs', page, status],
    queryFn: () => getJobs(page, status),
    placeholderData: previousData => previousData,
    retry: false,
  });

export const useJobsForQuery = (id: string, page: number) =>
  useQuery({
    queryKey: ['jobs', id, page],
    queryFn: () => getJobsFor(id, page),
    placeholderData: previousData => previousData,
    retry: false,
  });
