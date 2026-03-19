import { BaziResponse, getElementColor, getDirectionName } from '../types';
import { TrendingUp, MapPin, Star, Activity, Calendar } from 'lucide-react';

interface BaziChartProps {
  data: BaziResponse;
}

export default function BaziChart({ data }: BaziChartProps) {
  const pillars = ['year', 'month', 'day', 'hour'] as const;
  const pillarNames: Record<string, string> = {
    year: '年柱',
    month: '月柱',
    day: '日柱',
    hour: '時柱'
  };

  const getTenGodColor = (tenGod: string) => {
    if (tenGod.includes('比肩') || tenGod.includes('劫財')) return 'bg-blue-100 text-blue-800';
    if (tenGod.includes('食神') || tenGod.includes('傷官')) return 'bg-green-100 text-green-800';
    if (tenGod.includes('偏財') || tenGod.includes('正財')) return 'bg-yellow-100 text-yellow-800';
    if (tenGod.includes('七殺') || tenGod.includes('正官')) return 'bg-purple-100 text-purple-800';
    if (tenGod.includes('偏印') || tenGod.includes('正印')) return 'bg-red-100 text-red-800';
    return 'bg-gray-100 text-gray-800';
  };

  const getStrengthColor = (status: string) => {
    if (status.includes('強')) return 'text-red-600';
    if (status.includes('弱')) return 'text-blue-600';
    return 'text-gray-600';
  };

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="section-title flex items-center gap-2">
          <Calendar size={20} />
          四柱八字
        </h2>
        
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b-2 border-gray-200">
                <th className="py-2 px-3 text-left text-sm text-gray-500 w-24"></th>
                {pillars.map(pillar => (
                  <th key={pillar} className="py-2 px-3 text-center font-bold text-gray-700">
                    {pillarNames[pillar]}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="text-sm">
              <tr className="border-b border-gray-100">
                <td className="py-3 px-3 text-gray-500 font-medium">天干</td>
                {pillars.map(pillar => (
                  <td key={pillar} className="py-3 px-3 text-center">
                    <div className="stem-cell text-red-800">
                      {data.pillars[pillar].stem}
                    </div>
                  </td>
                ))}
              </tr>
              
              <tr className="border-b border-gray-100">
                <td className="py-2 px-3 text-gray-500">十神</td>
                {pillars.map(pillar => (
                  <td key={pillar} className="py-2 px-3 text-center">
                    <span className={`inline-block px-2 py-1 rounded text-xs font-medium ${getTenGodColor(data.pillars[pillar].ten_god_stem)}`}>
                      {data.pillars[pillar].ten_god_stem}
                    </span>
                  </td>
                ))}
              </tr>
              
              <tr className="border-b border-gray-100">
                <td className="py-3 px-3 text-gray-500 font-medium">地支</td>
                {pillars.map(pillar => (
                  <td key={pillar} className="py-3 px-3 text-center">
                    <div className="stem-cell text-red-800">
                      {data.pillars[pillar].branch}
                    </div>
                  </td>
                ))}
              </tr>
              
              <tr className="border-b border-gray-100">
                <td className="py-2 px-3 text-gray-500">藏干</td>
                {pillars.map(pillar => (
                  <td key={pillar} className="py-2 px-3 text-center text-xs">
                    <div className="space-y-1">
                      {data.pillars[pillar].hidden_stems.map((stem, idx) => (
                        <div key={idx} className="flex items-center justify-center gap-1">
                          <span className="text-gray-800">{stem}</span>
                          <span className="text-gray-500 text-xs">({data.pillars[pillar].ten_god_hidden[idx]})</span>
                        </div>
                      ))}
                    </div>
                  </td>
                ))}
              </tr>
              
              <tr className="border-b border-gray-100">
                <td className="py-2 px-3 text-gray-500">納音</td>
                {pillars.map(pillar => (
                  <td key={pillar} className="py-2 px-3 text-center text-xs text-gray-700">
                    {data.pillars[pillar].na_yin}
                  </td>
                ))}
              </tr>
              
              <tr className="border-b border-gray-100">
                <td className="py-2 px-3 text-gray-500">十二運</td>
                {pillars.map(pillar => (
                  <td key={pillar} className="py-2 px-3 text-center text-xs text-gray-700">
                    {data.pillars[pillar].life_stage}
                  </td>
                ))}
              </tr>
              
              <tr>
                <td className="py-2 px-3 text-gray-500">神煞</td>
                {pillars.map(pillar => (
                  <td key={pillar} className="py-2 px-3 text-center">
                    <div className="flex flex-wrap gap-1 justify-center">
                      {data.pillars[pillar].shen_sha.length > 0 ? (
                        data.pillars[pillar].shen_sha.map((sha, idx) => (
                          <span key={idx} className="shensha-tag">
                            {sha}
                          </span>
                        ))
                      ) : (
                        <span className="text-gray-400 text-xs">-</span>
                      )}
                    </div>
                  </td>
                ))}
              </tr>
            </tbody>
          </table>
        </div>

        <div className="mt-4 pt-4 border-t border-gray-200">
          <p className="text-sm text-gray-600">
            <span className="font-medium">日元：</span>
            <span className="text-xl font-bold text-red-800 mx-2">{data.day_stem}</span>
            <span className="text-gray-500">（命主天干）</span>
          </p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="section-title flex items-center gap-2">
          <Activity size={20} />
          身強身弱分析
        </h2>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-500 mb-1">狀態評估</p>
            <p className={`text-2xl font-bold ${getStrengthColor(data.strength.Status)}`}>
              {data.strength.Status}
            </p>
            <p className="text-sm text-gray-600 mt-1">
              分數: {data.strength.Score.toFixed(1)}%
            </p>
          </div>
          
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-500 mb-2">三得評估</p>
            <div className="space-y-1 text-sm">
              <div className="flex items-center gap-2">
                <span className={data.strength.IsDeLing ? 'text-green-600' : 'text-gray-400'}>
                  {data.strength.IsDeLing ? '✓' : '○'} 得令
                </span>
                <span className="text-xs text-gray-400">(月令)</span>
              </div>
              <div className="flex items-center gap-2">
                <span className={data.strength.IsDeDi ? 'text-green-600' : 'text-gray-400'}>
                  {data.strength.IsDeDi ? '✓' : '○'} 得地
                </span>
                <span className="text-xs text-gray-400">(地支)</span>
              </div>
              <div className="flex items-center gap-2">
                <span className={data.strength.IsDeZhu ? 'text-green-600' : 'text-gray-400'}>
                  {data.strength.IsDeZhu ? '✓' : '○'} 得助
                </span>
                <span className="text-xs text-gray-400">(天干)</span>
              </div>
            </div>
          </div>
          
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-500 mb-2">起運歲數</p>
            <p className="text-xl font-bold text-gray-800">
              {data.start_age_y} 歲 {data.start_age_m} 個月
            </p>
            <p className="text-xs text-gray-500 mt-1">開始進入第一柱大運</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="section-title flex items-center gap-2">
            <Star size={20} />
            喜忌五行
          </h2>
          
          <div className="space-y-3">
            <div>
              <p className="text-sm text-gray-500 mb-2">喜用神</p>
              <div className="flex flex-wrap gap-2">
                {(data.favorable_elements || []).map((el, idx) => (
                  <span key={idx} className={`element-badge ${getElementColor(el)}`}>
                    {el}
                  </span>
                ))}
                {(!data.favorable_elements || data.favorable_elements.length === 0) && (
                  <span className="text-gray-400 text-sm">暫無數據</span>
                )}
              </div>
            </div>
            
            <div className="pt-2 border-t border-gray-100">
              <p className="text-sm text-gray-500 mb-2">忌神</p>
              <div className="flex flex-wrap gap-2">
                {(data.unfavorable_elements || []).map((el, idx) => (
                  <span key={idx} className="element-badge bg-gray-200 text-gray-700">
                    {el}
                  </span>
                ))}
                {(!data.unfavorable_elements || data.unfavorable_elements.length === 0) && (
                  <span className="text-gray-400 text-sm">暫無數據</span>
                )}
              </div>
            </div>
          </div>
        </div>
        
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="section-title flex items-center gap-2">
            <MapPin size={20} />
            吉方位
          </h2>
          
          <div className="grid grid-cols-2 gap-3">
            <div className="bg-green-50 rounded-lg p-3">
              <p className="text-xs text-green-600 mb-1">財位</p>
              <p className="text-lg font-bold text-green-800">
                {data.directions ? getDirectionName(data.directions.wealth) : '暫無'}
              </p>
            </div>
            <div className="bg-blue-50 rounded-lg p-3">
              <p className="text-xs text-blue-600 mb-1">事業位</p>
              <p className="text-lg font-bold text-blue-800">
                {data.directions ? getDirectionName(data.directions.career) : '暫無'}
              </p>
            </div>
            <div className="bg-purple-50 rounded-lg p-3">
              <p className="text-xs text-purple-600 mb-1">文昌位</p>
              <p className="text-lg font-bold text-purple-800">
                {data.directions ? getDirectionName(data.directions.study) : '暫無'}
              </p>
            </div>
            <div className="bg-pink-50 rounded-lg p-3">
              <p className="text-xs text-pink-600 mb-1">桃花位</p>
              <p className="text-lg font-bold text-pink-800">
                {data.directions ? getDirectionName(data.directions.relationship) : '暫無'}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="section-title flex items-center gap-2">
          <TrendingUp size={20} />
          大運排列
        </h2>
        
        <div className="overflow-x-auto">
          <div className="flex gap-3 pb-2">
            {data.da_yun.map((dy, idx) => (
              <div key={idx} className="flex-shrink-0 bg-gray-50 rounded-lg p-3 min-w-[100px] text-center">
                <p className="text-xs text-gray-500 mb-1">{dy.start_age} 歲</p>
                <p className="text-lg font-bold text-gray-800">{dy.pillar}</p>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="section-title">
          流年斷語
        </h2>
        
        <div className="space-y-3">
          {data.advice.map((item, idx) => (
            <div key={idx} className={`p-3 rounded-lg border-l-4 ${
              item.Type === '吉' ? 'bg-green-50 border-green-500' :
              item.Type === '凶' ? 'bg-red-50 border-red-500' :
              'bg-gray-50 border-gray-400'
            }`}>
              <div className="flex items-center gap-2 mb-1">
                <span className={`text-xs font-bold px-2 py-0.5 rounded ${
                  item.Type === '吉' ? 'bg-green-200 text-green-800' :
                  item.Type === '凶' ? 'bg-red-200 text-red-800' :
                  'bg-gray-200 text-gray-700'
                }`}>
                  {item.Type}
                </span>
                <span className="font-medium text-gray-800">{item.Title}</span>
              </div>
              <p className="text-sm text-gray-600">{item.Content}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
