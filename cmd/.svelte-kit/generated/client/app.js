export { matchers } from './matchers.js';

export const nodes = [
	() => import('./nodes/0'),
	() => import('./nodes/1'),
	() => import('./nodes/2'),
	() => import('./nodes/3'),
	() => import('./nodes/4'),
	() => import('./nodes/5'),
	() => import('./nodes/6'),
	() => import('./nodes/7'),
	() => import('./nodes/8'),
	() => import('./nodes/9'),
	() => import('./nodes/10'),
	() => import('./nodes/11'),
	() => import('./nodes/12'),
	() => import('./nodes/13'),
	() => import('./nodes/14'),
	() => import('./nodes/15')
];

export const server_loads = [];

export const dictionary = {
		"/(login)": [4,[2]],
		"/(main)/backup": [5,[3]],
		"/(main)/export": [6,[3]],
		"/(main)/folders": [8,[3]],
		"/(main)/folder/[id]": [7,[3]],
		"/(main)/help": [9,[3]],
		"/(main)/import": [10,[3]],
		"/(main)/nodes": [11,[3]],
		"/(main)/recents": [12,[3]],
		"/(main)/rules": [13,[3]],
		"/(main)/search": [14,[3]],
		"/(main)/stats": [15,[3]]
	};

export const hooks = {
	handleError: (({ error }) => { console.error(error) }),
};

export { default as root } from '../root.svelte';