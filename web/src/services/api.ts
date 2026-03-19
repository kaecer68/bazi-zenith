import { BaziResponse } from '../types';

const API_BASE_URL = '/api/v1';

export async function fetchBaziChart(
  datetime: string,
  gender: string,
  targetYear: number = new Date().getFullYear()
): Promise<BaziResponse> {
  const response = await fetch(`${API_BASE_URL}/chart`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      datetime,
      gender,
      target_year: targetYear
    }),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to fetch Bazi chart');
  }

  return response.json();
}

export async function checkHealth(): Promise<boolean> {
  try {
    const response = await fetch('/health');
    return response.ok;
  } catch {
    return false;
  }
}
