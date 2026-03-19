import { useState, useEffect } from 'react';
import { BirthData } from '../types';
import { getAllProfiles, deleteProfile } from '../utils/storage';
import { User, Trash2, Plus, Search } from 'lucide-react';

interface ProfileListProps {
  onSelect: (profile: BirthData | null) => void;
  selectedId?: string;
}

export default function ProfileList({ onSelect, selectedId }: ProfileListProps) {
  const [profiles, setProfiles] = useState<BirthData[]>([]);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    setProfiles(getAllProfiles());
  }, []);

  const handleDelete = (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    if (confirm('確定要刪除此生辰資料嗎？')) {
      deleteProfile(id);
      setProfiles(getAllProfiles());
      if (selectedId === id) {
        onSelect(null);
      }
    }
  };

  const handleNew = () => {
    onSelect(null);
  };

  const filteredProfiles = profiles.filter(p => 
    p.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const formatDate = (profile: BirthData) => {
    const cal = profile.calendarType === 'solar' ? '陽曆' : '農曆';
    return `${cal} ${profile.year}-${String(profile.month).padStart(2, '0')}-${String(profile.day).padStart(2, '0')}`;
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-bold text-gray-800">生辰資料管理</h2>
        <button
          onClick={handleNew}
          className="btn-primary flex items-center gap-1 text-sm"
        >
          <Plus size={16} />
          新增
        </button>
      </div>

      <div className="relative mb-4">
        <Search className="absolute left-3 top-2.5 text-gray-400" size={18} />
        <input
          type="text"
          placeholder="搜尋姓名..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="input-field pl-10"
        />
      </div>

      <div className="space-y-2 max-h-96 overflow-y-auto">
        {filteredProfiles.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <User className="mx-auto mb-2" size={32} />
            <p className="text-sm">暫無生辰資料</p>
            <p className="text-xs text-gray-400 mt-1">點擊「新增」建立第一筆資料</p>
          </div>
        ) : (
          filteredProfiles.map(profile => (
            <div
              key={profile.id}
              onClick={() => onSelect(profile)}
              className={`p-3 rounded-lg cursor-pointer transition-colors ${
                selectedId === profile.id 
                  ? 'bg-red-50 border-2 border-red-200' 
                  : 'bg-gray-50 hover:bg-gray-100 border-2 border-transparent'
              }`}
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-800">{profile.name}</p>
                  <p className="text-xs text-gray-500 mt-1">
                    {formatDate(profile)} {String(profile.hour).padStart(2, '0')}:{String(profile.minute).padStart(2, '0')}
                  </p>
                  <p className="text-xs text-gray-400">
                    {profile.gender === 'male' ? '乾造 (男)' : '坤造 (女)'}
                  </p>
                </div>
                <button
                  onClick={(e) => handleDelete(e, profile.id)}
                  className="p-2 text-gray-400 hover:text-red-600 transition-colors"
                  title="刪除"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
