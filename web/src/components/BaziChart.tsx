import { useEffect, useMemo, useState } from 'react';
import { BaziResponse, ChartQueryContext, getDirectionName, getElementColor } from '../types';
import { TrendingUp, MapPin, Star, Activity, Calendar } from 'lucide-react';

interface BaziChartProps {
  data: BaziResponse;
  queryContext: ChartQueryContext | null;
  onTargetYearChange: (year: number) => void | Promise<void>;
  isSwitchingYear: boolean;
}

export default function BaziChart({ data, queryContext, onTargetYearChange, isSwitchingYear }: BaziChartProps) {
  const pillars = ['year', 'month', 'day', 'hour'] as const;
  const pillarNames: Record<string, string> = {
    year: '年柱',
    month: '月柱',
    day: '日柱',
    hour: '時柱'
  };

  const getYearPillar = (year: number) => {
    const stems = ['甲', '乙', '丙', '丁', '戊', '己', '庚', '辛', '壬', '癸'];
    const branches = ['子', '丑', '寅', '卯', '辰', '巳', '午', '未', '申', '酉', '戌', '亥'];
    let idx = (year - 4) % 60;
    if (idx < 0) idx += 60;
    return `${stems[idx % 10]}${branches[idx % 12]}`;
  };

  const detail = data.detail_chart;
  const dayunBoard = detail?.dayun_board?.length
    ? detail.dayun_board
    : data.da_yun.map((item, idx) => ({
        index: idx + 1,
        year: 0,
        start_age: item.start_age,
        start_year: undefined,
        pillar: item.pillar,
        ten_god_stem: '',
        ten_god_branch: '',
      }));
  const liunianBoard = detail?.liunian_board || [];
  const targetYear = queryContext?.targetYear || liunianBoard[0]?.year || new Date().getFullYear();
  const defaultDayunIndex = useMemo(() => {
    if (dayunBoard.length === 0) return 0;
    const idx = dayunBoard.findIndex((item, i) => {
      const nextStart = dayunBoard[i + 1]?.start_year;
      const inLowerBound = item.start_year === undefined || targetYear >= item.start_year;
      const inUpperBound = nextStart === undefined || targetYear < nextStart;
      return inLowerBound && inUpperBound;
    });
    return idx >= 0 ? idx : 0;
  }, [dayunBoard, targetYear]);
  const [selectedDayunIndex, setSelectedDayunIndex] = useState(defaultDayunIndex);

  useEffect(() => {
    setSelectedDayunIndex(defaultDayunIndex);
  }, [defaultDayunIndex]);

  const selectedDayun = dayunBoard[selectedDayunIndex];
  const yearCandidates = useMemo(() => {
    if (selectedDayun?.start_year) {
      return Array.from({ length: 10 }, (_, i) => selectedDayun.start_year! + i);
    }
    if (liunianBoard.length > 0) {
      return liunianBoard.map(item => item.year);
    }
    return Array.from({ length: 10 }, (_, i) => targetYear - 4 + i);
  }, [selectedDayun, liunianBoard, targetYear]);

  const selectedLiunian = useMemo(() => {
    if (liunianBoard.length > 0) {
      return liunianBoard.find(item => item.year === targetYear) || liunianBoard[0];
    }
    return {
      year: targetYear,
      pillar: getYearPillar(targetYear),
      ten_god_stem: '',
      ten_god_branch: '',
    };
  }, [liunianBoard, targetYear]);

  const activeLiunianPillar = selectedLiunian?.pillar || '';
  const activeDayunPillar = selectedDayun?.pillar || '';

  const splitPillar = (pillar: string): [string, string] => {
    const chars = Array.from(pillar || '');
    return [chars[0] || '', chars[1] || ''];
  };

  const [dayunStem, dayunBranch] = splitPillar(activeDayunPillar);
  const [liunianStem, liunianBranch] = splitPillar(activeLiunianPillar);

  const stemHePairs = new Set(['甲己', '己甲', '乙庚', '庚乙', '丙辛', '辛丙', '丁壬', '壬丁', '戊癸', '癸戊']);
  const stemChongPairs = new Set(['甲庚', '庚甲', '乙辛', '辛乙', '丙壬', '壬丙', '丁癸', '癸丁']);
  const branchHePairs = new Set(['子丑', '丑子', '寅亥', '亥寅', '卯戌', '戌卯', '辰酉', '酉辰', '巳申', '申巳', '午未', '未午']);
  const branchChongPairs = new Set(['子午', '午子', '丑未', '未丑', '寅申', '申寅', '卯酉', '酉卯', '辰戌', '戌辰', '巳亥', '亥巳']);
  const branchHaiPairs = new Set(['子未', '未子', '丑午', '午丑', '寅巳', '巳寅', '卯辰', '辰卯', '申亥', '亥申', '酉戌', '戌酉']);
  const branchXingPairs = new Set([
    '子卯', '卯子',
    '丑未', '未丑', '丑戌', '戌丑', '未戌', '戌未',
    '寅巳', '巳寅', '寅申', '申寅', '巳申', '申巳',
    '辰辰', '午午', '酉酉', '亥亥',
  ]);

  const isBranchXing = (a: string, b: string) => {
    const pair = `${a}${b}`;
    return branchXingPairs.has(pair);
  };

  const stemRel = (a: string, b: string) => {
    const key = `${a}${b}`;
    const tags: string[] = [];
    if (stemHePairs.has(key)) tags.push('合');
    if (stemChongPairs.has(key)) tags.push('沖');
    return tags;
  };

  const branchRel = (a: string, b: string) => {
    const key = `${a}${b}`;
    const tags: string[] = [];
    if (branchHePairs.has(key)) tags.push('合');
    if (branchChongPairs.has(key)) tags.push('沖');
    if (isBranchXing(a, b)) tags.push('刑');
    if (branchHaiPairs.has(key)) tags.push('害');
    return tags;
  };

  const interactions = useMemo(() => {
    const nodes = [
      {
        label: '年柱',
        stem: detail?.natal.tian_gan.year || data.pillars.year.stem,
        branch: detail?.natal.di_zhi.year || data.pillars.year.branch,
      },
      {
        label: '月柱',
        stem: detail?.natal.tian_gan.month || data.pillars.month.stem,
        branch: detail?.natal.di_zhi.month || data.pillars.month.branch,
      },
      {
        label: '日柱',
        stem: detail?.natal.tian_gan.day || data.pillars.day.stem,
        branch: detail?.natal.di_zhi.day || data.pillars.day.branch,
      },
      {
        label: '時柱',
        stem: detail?.natal.tian_gan.hour || data.pillars.hour.stem,
        branch: detail?.natal.di_zhi.hour || data.pillars.hour.branch,
      },
      ...(dayunStem || dayunBranch ? [{ label: '大運', stem: dayunStem, branch: dayunBranch }] : []),
      ...(liunianStem || liunianBranch ? [{ label: '流年', stem: liunianStem, branch: liunianBranch }] : []),
    ];

    const stemNotes: string[] = [];
    const branchNotes: string[] = [];
    const stemSeen = new Set<string>();
    const branchSeen = new Set<string>();

    for (let i = 0; i < nodes.length; i += 1) {
      for (let j = i + 1; j < nodes.length; j += 1) {
        const a = nodes[i];
        const b = nodes[j];

        if (a.stem && b.stem) {
          const tags = stemRel(a.stem, b.stem);
          if (tags.length > 0) {
            const key = `${a.label}-${b.label}-${a.stem}-${b.stem}-${tags.join('/')}`;
            if (!stemSeen.has(key)) {
              stemSeen.add(key);
              stemNotes.push(`${a.label}${a.stem} × ${b.label}${b.stem}：${tags.join('/')}`);
            }
          }
        }

        if (a.branch && b.branch) {
          const tags = branchRel(a.branch, b.branch);
          if (tags.length > 0) {
            const key = `${a.label}-${b.label}-${a.branch}-${b.branch}-${tags.join('/')}`;
            if (!branchSeen.has(key)) {
              branchSeen.add(key);
              branchNotes.push(`${a.label}${a.branch} × ${b.label}${b.branch}：${tags.join('/')}`);
            }
          }
        }
      }
    }

    return { stem: stemNotes, branch: branchNotes };
  }, [detail, data.pillars, dayunStem, dayunBranch, liunianStem, liunianBranch]);

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

  const strengthStatus = data.strength.status || data.strength.Status || '未定';
  const strengthScore = data.strength.score ?? data.strength.Score ?? 0;
  const isDeLing = data.strength.is_de_ling ?? data.strength.IsDeLing;
  const isDeDi = data.strength.is_de_di ?? data.strength.IsDeDi;
  const isDeZhu = data.strength.is_de_zhu ?? data.strength.IsDeZhu;

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
        <h2 className="section-title">沖合刑害即時分析</h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="border border-red-100 bg-red-50 rounded-lg p-3">
              <p className="text-sm font-medium text-red-900 mb-2">天干（四柱＋大運＋流年）</p>
              <div className="space-y-1 text-sm text-gray-700">
                {interactions.stem.length > 0 ? interactions.stem.map((msg, idx) => <p key={idx}>{msg}</p>) : <p>目前無明顯沖合。</p>}
              </div>
            </div>

            <div className="border border-blue-100 bg-blue-50 rounded-lg p-3">
              <p className="text-sm font-medium text-blue-900 mb-2">地支（四柱＋大運＋流年）</p>
              <div className="space-y-1 text-sm text-gray-700">
                {interactions.branch.length > 0 ? interactions.branch.map((msg, idx) => <p key={idx}>{msg}</p>) : <p>目前無明顯沖合刑害。</p>}
              </div>
            </div>
          </div>

        <div className="mt-4 text-sm text-gray-600 bg-gray-50 rounded p-3 border border-gray-200">
          <p className="font-medium mb-1">細盤六柱提示</p>
          <p>{detail?.prompts.tiangan || '—'}</p>
          <p>{detail?.prompts.dizhi || '—'}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="section-title flex items-center gap-2">
          <TrendingUp size={20} />
          大運 / 流年互動切換
        </h2>

          <div className="space-y-3 border border-orange-100 bg-orange-50 rounded-lg p-3">
            <div>
              <p className="text-sm text-orange-900 font-medium mb-2">大運切換</p>
              <div className="flex gap-2 overflow-x-auto pb-1">
                {dayunBoard.map((item, idx) => (
                  <button
                    key={`${item.pillar}-${item.start_year || idx}`}
                    type="button"
                    onClick={() => setSelectedDayunIndex(idx)}
                    className={`px-3 py-2 rounded border text-sm whitespace-nowrap ${
                      idx === selectedDayunIndex
                        ? 'bg-orange-600 text-white border-orange-600'
                        : 'bg-white text-gray-700 border-orange-200 hover:bg-orange-100'
                    }`}
                  >
                    {item.start_age || 0}歲 · {item.pillar}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <p className="text-sm text-orange-900 font-medium mb-2">流年切換（即時重算）</p>
              <div className="flex gap-2 overflow-x-auto pb-1">
                {yearCandidates.map((year) => (
                  <button
                    key={year}
                    type="button"
                    disabled={isSwitchingYear}
                    onClick={() => onTargetYearChange(year)}
                    className={`px-3 py-1.5 rounded border text-sm whitespace-nowrap ${
                      year === targetYear
                        ? 'bg-red-700 text-white border-red-700'
                        : 'bg-white text-gray-700 border-red-200 hover:bg-red-50'
                    } ${isSwitchingYear ? 'opacity-60 cursor-not-allowed' : ''}`}
                  >
                    {year}
                  </button>
                ))}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
              <div className="bg-white rounded border border-gray-200 p-2">
                <p className="text-gray-500">當前大運</p>
                <p className="font-semibold text-gray-800">
                  {activeDayunPillar || '—'}
                  {selectedDayun?.start_year ? `（${selectedDayun.start_year}-${selectedDayun.start_year + 9}）` : ''}
                </p>
              </div>
              <div className="bg-white rounded border border-gray-200 p-2">
                <p className="text-gray-500">當前流年</p>
                <p className="font-semibold text-gray-800">{targetYear} · {activeLiunianPillar || '—'}</p>
              </div>
            </div>
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
            <p className={`text-2xl font-bold ${getStrengthColor(strengthStatus)}`}>
              {strengthStatus}
            </p>
            <p className="text-sm text-gray-600 mt-1">
              分數: {strengthScore.toFixed(1)}%
            </p>
          </div>
          
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-500 mb-2">三得評估</p>
            <div className="space-y-1 text-sm">
              <div className="flex items-center gap-2">
                <span className={isDeLing ? 'text-green-600' : 'text-gray-400'}>
                  {isDeLing ? '✓' : '○'} 得令
                </span>
                <span className="text-xs text-gray-400">(月令)</span>
              </div>
              <div className="flex items-center gap-2">
                <span className={isDeDi ? 'text-green-600' : 'text-gray-400'}>
                  {isDeDi ? '✓' : '○'} 得地
                </span>
                <span className="text-xs text-gray-400">(地支)</span>
              </div>
              <div className="flex items-center gap-2">
                <span className={isDeZhu ? 'text-green-600' : 'text-gray-400'}>
                  {isDeZhu ? '✓' : '○'} 得助
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
        <h2 className="section-title">
          流年斷語
        </h2>
        
        <div className="space-y-3">
          {data.advice.map((item, idx) => {
            const adviceType = item.type || item.Type || '平';
            const adviceTitle = item.title || item.Title || '提示';
            const adviceContent = item.content || item.Content || '';
            return (
            <div key={idx} className={`p-3 rounded-lg border-l-4 ${
              adviceType === '吉' ? 'bg-green-50 border-green-500' :
              adviceType === '凶' ? 'bg-red-50 border-red-500' :
              'bg-gray-50 border-gray-400'
            }`}>
              <div className="flex items-center gap-2 mb-1">
                <span className={`text-xs font-bold px-2 py-0.5 rounded ${
                  adviceType === '吉' ? 'bg-green-200 text-green-800' :
                  adviceType === '凶' ? 'bg-red-200 text-red-800' :
                  'bg-gray-200 text-gray-700'
                }`}>
                  {adviceType}
                </span>
                <span className="font-medium text-gray-800">{adviceTitle}</span>
              </div>
              <p className="text-sm text-gray-600">{adviceContent}</p>
            </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
