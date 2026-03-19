import { useState } from 'react';
import ProfileList from './components/ProfileList';
import BirthForm from './components/BirthForm';
import BaziChart from './components/BaziChart';
import { BirthData, BaziResponse } from './types';
import { testChartData } from './testData';

function App() {
  const [selectedProfile, setSelectedProfile] = useState<BirthData | null>(null);
  const [chartData, setChartData] = useState<BaziResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleProfileSelect = (profile: BirthData | null) => {
    setSelectedProfile(profile);
    setChartData(null);
    setError(null);
  };

  const handleChartGenerated = (data: BaziResponse) => {
    console.log('Chart data received:', data);
    setChartData(data);
    setIsLoading(false);
  };

  const handleLoading = (loading: boolean) => {
    setIsLoading(loading);
  };

  const handleError = (err: string | null) => {
    setError(err);
    setIsLoading(false);
  };

  const loadTestData = () => {
    console.log('Loading test data:', testChartData);
    setChartData(testChartData as unknown as BaziResponse);
    setError(null);
  };

  return (
    <div className="min-h-screen bg-stone-50">
      <header className="bg-red-900 text-white py-4 shadow-lg">
        <div className="container mx-auto px-4">
          <h1 className="text-2xl font-bold tracking-wide">八字命盤管理系統</h1>
          <p className="text-red-200 text-sm mt-1">Bazi Zenith Chart Management</p>
        </div>
      </header>

      <main className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <div className="lg:col-span-4 space-y-6">
            <ProfileList 
              onSelect={handleProfileSelect}
              selectedId={selectedProfile?.id}
            />
            
            <div className="bg-blue-50 rounded-lg p-4 border border-blue-200">
              <p className="text-sm text-blue-800 mb-2">開發測試</p>
              <button
                onClick={loadTestData}
                className="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded text-sm"
              >
                載入測試數據
              </button>
            </div>
          </div>

          <div className="lg:col-span-8 space-y-6">
            <BirthForm 
              selectedProfile={selectedProfile}
              onChartGenerated={handleChartGenerated}
              onLoading={handleLoading}
              onError={handleError}
            />

            {isLoading && (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-red-800"></div>
                <span className="ml-3 text-gray-600">計算八字命盤中...</span>
              </div>
            )}

            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-800">
                <p className="font-medium">錯誤</p>
                <p className="text-sm">{error}</p>
              </div>
            )}

            {chartData && !isLoading && (
              <BaziChart data={chartData} />
            )}
          </div>
        </div>
      </main>

      <footer className="bg-gray-800 text-gray-300 py-4 mt-12">
        <div className="container mx-auto px-4 text-center text-sm">
          <p>八字命盤引擎 Bazi-Zenith | 基於天文級黃曆計算</p>
        </div>
      </footer>
    </div>
  );
}

export default App;
