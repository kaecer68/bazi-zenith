import { useState, useEffect } from 'react';
import { BirthData, BaziResponse, convertToSolarDate } from '../types';
import { saveProfile, generateId } from '../utils/storage';
import { fetchBaziChart } from '../services/api';
import { Calendar, Clock, User, Save, Calculator } from 'lucide-react';

interface BirthFormProps {
  selectedProfile: BirthData | null;
  onChartGenerated: (data: BaziResponse) => void;
  onLoading: (loading: boolean) => void;
  onError: (error: string | null) => void;
}

export default function BirthForm({ 
  selectedProfile, 
  onChartGenerated, 
  onLoading, 
  onError 
}: BirthFormProps) {
  const [formData, setFormData] = useState<Partial<BirthData>>({
    name: '',
    gender: 'male',
    calendarType: 'solar',
    year: new Date().getFullYear() - 30,
    month: 1,
    day: 1,
    hour: 12,
    minute: 0,
    note: ''
  });
  const [targetYear, setTargetYear] = useState(new Date().getFullYear());

  useEffect(() => {
    if (selectedProfile) {
      setFormData(selectedProfile);
    }
  }, [selectedProfile]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.name || !formData.year || !formData.month || !formData.day) {
      onError('請填寫完整的生辰資料');
      return;
    }

    onLoading(true);
    onError(null);

    try {
      const birthData = formData as BirthData;
      const solarDate = convertToSolarDate(birthData);
      const hourStr = String(birthData.hour).padStart(2, '0');
      const minuteStr = String(birthData.minute).padStart(2, '0');
      const datetime = `${solarDate.year}-${String(solarDate.month).padStart(2, '0')}-${String(solarDate.day).padStart(2, '0')} ${hourStr}:${minuteStr}`;
      
      const chartData = await fetchBaziChart(datetime, birthData.gender, targetYear);
      onChartGenerated(chartData);

      const profile: BirthData = {
        ...birthData,
        id: selectedProfile?.id || generateId(),
        createdAt: selectedProfile?.createdAt || Date.now()
      };
      saveProfile(profile);
    } catch (err) {
      onError(err instanceof Error ? err.message : '計算失敗');
    } finally {
      onLoading(false);
    }
  };

  const handleChange = (field: keyof BirthData, value: string | number) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="section-title flex items-center gap-2">
        <Calendar size={20} />
        生辰輸入
      </h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              <User size={14} className="inline mr-1" />
              姓名
            </label>
            <input
              type="text"
              value={formData.name || ''}
              onChange={(e) => handleChange('name', e.target.value)}
              placeholder="輸入姓名"
              className="input-field"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">性別</label>
            <div className="flex gap-4">
              <label className="flex items-center">
                <input
                  type="radio"
                  value="male"
                  checked={formData.gender === 'male'}
                  onChange={(e) => handleChange('gender', e.target.value)}
                  className="mr-2"
                />
                <span>乾造 (男)</span>
              </label>
              <label className="flex items-center">
                <input
                  type="radio"
                  value="female"
                  checked={formData.gender === 'female'}
                  onChange={(e) => handleChange('gender', e.target.value)}
                  className="mr-2"
                />
                <span>坤造 (女)</span>
              </label>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">曆法類型</label>
            <select
              value={formData.calendarType || 'solar'}
              onChange={(e) => handleChange('calendarType', e.target.value)}
              className="input-field"
            >
              <option value="solar">陽曆 (西曆)</option>
              <option value="lunar">農曆 (陰曆)</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              <Clock size={14} className="inline mr-1" />
              流年目標年份
            </label>
            <input
              type="number"
              value={targetYear}
              onChange={(e) => setTargetYear(parseInt(e.target.value))}
              min="1900"
              max="2100"
              className="input-field"
            />
          </div>
        </div>

        <div className="border-t border-gray-200 pt-4">
          <h3 className="text-sm font-medium text-gray-700 mb-3">
            {formData.calendarType === 'solar' ? '陽曆' : '農曆'}日期
          </h3>
          
          <div className="grid grid-cols-3 gap-3">
            <div>
              <label className="block text-xs text-gray-500 mb-1">年</label>
              <input
                type="number"
                value={formData.year || ''}
                onChange={(e) => handleChange('year', parseInt(e.target.value))}
                min="1900"
                max="2100"
                className="input-field"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">月</label>
              <input
                type="number"
                value={formData.month || ''}
                onChange={(e) => handleChange('month', parseInt(e.target.value))}
                min="1"
                max="12"
                className="input-field"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">日</label>
              <input
                type="number"
                value={formData.day || ''}
                onChange={(e) => handleChange('day', parseInt(e.target.value))}
                min="1"
                max="31"
                className="input-field"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3 mt-3">
            <div>
              <label className="block text-xs text-gray-500 mb-1">時</label>
              <input
                type="number"
                value={formData.hour || ''}
                onChange={(e) => handleChange('hour', parseInt(e.target.value))}
                min="0"
                max="23"
                className="input-field"
              />
            </div>
            <div>
              <label className="block text-xs text-gray-500 mb-1">分</label>
              <input
                type="number"
                value={formData.minute || ''}
                onChange={(e) => handleChange('minute', parseInt(e.target.value))}
                min="0"
                max="59"
                className="input-field"
              />
            </div>
          </div>
        </div>

        <div className="flex gap-3 pt-4">
          <button
            type="submit"
            className="btn-primary flex-1 flex items-center justify-center gap-2"
          >
            <Calculator size={18} />
            排盤計算
          </button>
          <button
            type="button"
            onClick={() => {
              const profile: BirthData = {
                ...(formData as BirthData),
                id: selectedProfile?.id || generateId(),
                createdAt: selectedProfile?.createdAt || Date.now()
              };
              saveProfile(profile);
              alert('資料已儲存');
            }}
            className="btn-secondary flex items-center gap-2"
          >
            <Save size={18} />
            儲存
          </button>
        </div>
      </form>
    </div>
  );
}
