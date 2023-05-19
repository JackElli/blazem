export const manifest = {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["favicon.png"]),
	mimeTypes: {".png":"image/png"},
	_: {
		client: {"start":{"file":"_app/immutable/entry/start.92310e54.js","imports":["_app/immutable/entry/start.92310e54.js","_app/immutable/chunks/index.880e7ca1.js","_app/immutable/chunks/singletons.d83749f5.js"],"stylesheets":[],"fonts":[]},"app":{"file":"_app/immutable/entry/app.b367bfe3.js","imports":["_app/immutable/entry/app.b367bfe3.js","_app/immutable/chunks/index.880e7ca1.js"],"stylesheets":[],"fonts":[]}},
		nodes: [
			() => import('./nodes/0.js'),
			() => import('./nodes/1.js'),
			() => import('./nodes/2.js'),
			() => import('./nodes/3.js'),
			() => import('./nodes/4.js'),
			() => import('./nodes/5.js'),
			() => import('./nodes/6.js'),
			() => import('./nodes/7.js'),
			() => import('./nodes/8.js'),
			() => import('./nodes/9.js'),
			() => import('./nodes/10.js'),
			() => import('./nodes/11.js'),
			() => import('./nodes/12.js'),
			() => import('./nodes/13.js'),
			() => import('./nodes/14.js'),
			() => import('./nodes/15.js'),
			() => import('./nodes/16.js')
		],
		routes: [
			{
				id: "/(login)",
				pattern: /^\/?$/,
				params: [],
				page: { layouts: [0,2], errors: [1,,], leaf: 5 },
				endpoint: null
			},
			{
				id: "/(main)/backup",
				pattern: /^\/backup\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 6 },
				endpoint: null
			},
			{
				id: "/(main)/export",
				pattern: /^\/export\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 7 },
				endpoint: null
			},
			{
				id: "/(main)/folders",
				pattern: /^\/folders\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 9 },
				endpoint: null
			},
			{
				id: "/(main)/folder/[id]",
				pattern: /^\/folder\/([^/]+?)\/?$/,
				params: [{"name":"id","optional":false,"rest":false,"chained":false}],
				page: { layouts: [0,3,4], errors: [1,,,], leaf: 8 },
				endpoint: null
			},
			{
				id: "/(main)/help",
				pattern: /^\/help\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 10 },
				endpoint: null
			},
			{
				id: "/(main)/import",
				pattern: /^\/import\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 11 },
				endpoint: null
			},
			{
				id: "/(main)/nodes",
				pattern: /^\/nodes\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 12 },
				endpoint: null
			},
			{
				id: "/(main)/recents",
				pattern: /^\/recents\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 13 },
				endpoint: null
			},
			{
				id: "/(main)/rules",
				pattern: /^\/rules\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 14 },
				endpoint: null
			},
			{
				id: "/(main)/search",
				pattern: /^\/search\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 15 },
				endpoint: null
			},
			{
				id: "/(main)/stats",
				pattern: /^\/stats\/?$/,
				params: [],
				page: { layouts: [0,3], errors: [1,,], leaf: 16 },
				endpoint: null
			}
		],
		matchers: async () => {
			
			return {  };
		}
	}
};
