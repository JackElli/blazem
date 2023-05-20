const manifest = {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["favicon.png"]),
	mimeTypes: {".png":"image/png"},
	_: {
		client: {"start":{"file":"_app/immutable/entry/start.92310e54.js","imports":["_app/immutable/entry/start.92310e54.js","_app/immutable/chunks/index.880e7ca1.js","_app/immutable/chunks/singletons.d83749f5.js"],"stylesheets":[],"fonts":[]},"app":{"file":"_app/immutable/entry/app.b367bfe3.js","imports":["_app/immutable/entry/app.b367bfe3.js","_app/immutable/chunks/index.880e7ca1.js"],"stylesheets":[],"fonts":[]}},
		nodes: [
			() => import('./chunks/0-149650e6.js'),
			() => import('./chunks/1-1c17bc9c.js'),
			() => import('./chunks/2-5e3943fd.js'),
			() => import('./chunks/3-6b53724c.js'),
			() => import('./chunks/4-b12242a5.js'),
			() => import('./chunks/5-7814d26d.js'),
			() => import('./chunks/6-11ce5d0c.js'),
			() => import('./chunks/7-16b6972e.js'),
			() => import('./chunks/8-3b3b0e5d.js'),
			() => import('./chunks/9-ec1e95ad.js'),
			() => import('./chunks/10-042b65ab.js'),
			() => import('./chunks/11-88507c26.js'),
			() => import('./chunks/12-ca2970a9.js'),
			() => import('./chunks/13-96136344.js'),
			() => import('./chunks/14-92541804.js'),
			() => import('./chunks/15-7957eb75.js'),
			() => import('./chunks/16-e4e47736.js')
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

const prerendered = new Set([]);

export { manifest, prerendered };
//# sourceMappingURL=manifest.js.map
