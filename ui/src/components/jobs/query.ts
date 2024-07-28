import axios from 'axios';

import { useQuery } from '@tanstack/react-query';

import { JobsResponse, Stats } from './types';

export const getJobsFor = async (id: string, page: number) => {
  const response = await axios.get(`/api/minion/jobs/?page=${page}&client=${id}`);
  return response.data.jobs as JobsResponse;
};

export const getJobs = async (page: number, status: string, client: string) => {
  const response = await axios.get(`/api/minion/jobs?limit=100&page=${page}&status=${status}&client=${client}`);
  return response.data as JobsResponse;
};
export const getJobStats = async () => {
  const response = await axios.get(`/api/minion/jobs?limit=1`);
  return response.data.stats as Stats;
};

export const queueJob = async (name: string, client: string) => {
  const response = await axios.post(`/api/minion/jobs?job=${name}&client=${client}`);
  return response.data;
};

export const deleteJob = async (id: string, hard: boolean) => {
  const response = await axios.delete(`/api/minion/jobs/${id}?hard=${hard}`);
  return response.data;
};

export const patchJob = async (id: string) => {
  const response = await axios.patch(`/api/minion/jobs/${id}`, {});
  return response.data;
};

export const useJobsStatsQuery = () =>
  useQuery({
    queryKey: ['jobs', 'stats'],
    queryFn: () => getJobStats(),
    placeholderData: previousData => previousData,
    retry: false,
  });

export const useJobsQuery = (page: number, status: string, client: string) =>
  useQuery({
    queryKey: ['jobs', page, status],
    queryFn: () => getJobs(page, status, client),
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
