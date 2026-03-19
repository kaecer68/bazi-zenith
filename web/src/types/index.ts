import { Solar, Lunar } from 'lunar-javascript';

export interface BirthData {
  id: string;
  name: string;
  gender: 'male' | 'female';
  calendarType: 'solar' | 'lunar';
  year: number;
  month: number;
  day: number;
  hour: number;
  minute: number;
  note?: string;
  createdAt: number;
}

export interface PillarData {
  stem: string;
  branch: string;
  ten_god_stem: string;
  hidden_stems: string[];
  ten_god_hidden: string[];
  na_yin: string;
  life_stage: string;
  shen_sha: string[];
}

export interface DaYunData {
  pillar: string;
  start_age: number;
}

export interface StrengthAnalysis {
  Score: number;
  Status: string;
  IsDeLing: boolean;
  IsDeDi: boolean;
  IsDeZhu: boolean;
  Percentage: number;
}

export interface Interpretation {
  Title: string;
  Content: string;
  Type: string;
}

export interface BaziDirections {
  wealth: string;
  career: string;
  study: string;
  relationship: string;
}

export interface BaziResponse {
  gender: string;
  day_stem: string;
  pillars: {
    year: PillarData;
    month: PillarData;
    day: PillarData;
    hour: PillarData;
  };
  da_yun: DaYunData[];
  start_age_y: number;
  start_age_m: number;
  strength: StrengthAnalysis;
  advice: Interpretation[];
  favorable_elements: string[];
  unfavorable_elements: string[];
  directions: BaziDirections;
}

export function convertToSolarDate(birthData: BirthData): { year: number; month: number; day: number } {
  if (birthData.calendarType === 'solar') {
    return {
      year: birthData.year,
      month: birthData.month,
      day: birthData.day
    };
  }
  
  const lunar = Lunar.fromYmd(birthData.year, birthData.month, birthData.day);
  const solar = lunar.getSolar();
  return {
    year: solar.getYear(),
    month: solar.getMonth(),
    day: solar.getDay()
  };
}

export function convertToLunarDate(year: number, month: number, day: number): { year: number; month: number; day: number } {
  const solar = Solar.fromYmd(year, month, day);
  const lunar = solar.getLunar();
  return {
    year: lunar.getYear(),
    month: lunar.getMonth(),
    day: lunar.getDay()
  };
}

export function formatDateTime(birthData: BirthData): string {
  const date = convertToSolarDate(birthData);
  const hourStr = birthData.hour.toString().padStart(2, '0');
  const minuteStr = birthData.minute.toString().padStart(2, '0');
  return `${date.year}-${date.month.toString().padStart(2, '0')}-${date.day.toString().padStart(2, '0')} ${hourStr}:${minuteStr}`;
}

export function getElementColor(element: string): string {
  const colorMap: Record<string, string> = {
    '木': 'bg-green-600 text-white',
    '火': 'bg-red-600 text-white',
    '土': 'bg-yellow-500 text-black',
    '金': 'bg-gray-400 text-black',
    '水': 'bg-blue-600 text-white',
  };
  return colorMap[element] || 'bg-gray-200 text-black';
}

export function getDirectionName(direction: string): string {
  const directionMap: Record<string, string> = {
    '東': '東方',
    '南': '南方',
    '西': '西方',
    '北': '北方',
    '東南': '東南方',
    '西南': '西南方',
    '東北': '東北方',
    '西北': '西北方',
  };
  return directionMap[direction] || direction;
}
