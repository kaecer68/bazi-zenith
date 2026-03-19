import { BirthData } from '../types';

const STORAGE_KEY = 'bazi-zenith-profiles';

export function getAllProfiles(): BirthData[] {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (!stored) return [];
  try {
    return JSON.parse(stored);
  } catch {
    return [];
  }
}

export function saveProfile(profile: BirthData): void {
  const profiles = getAllProfiles();
  const existingIndex = profiles.findIndex(p => p.id === profile.id);
  
  if (existingIndex >= 0) {
    profiles[existingIndex] = profile;
  } else {
    profiles.push(profile);
  }
  
  localStorage.setItem(STORAGE_KEY, JSON.stringify(profiles));
}

export function deleteProfile(id: string): void {
  const profiles = getAllProfiles().filter(p => p.id !== id);
  localStorage.setItem(STORAGE_KEY, JSON.stringify(profiles));
}

export function getProfileById(id: string): BirthData | undefined {
  return getAllProfiles().find(p => p.id === id);
}

export function generateId(): string {
  return Date.now().toString(36) + Math.random().toString(36).substr(2);
}
