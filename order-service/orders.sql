

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.orders_id_seq OWNER TO postgres2;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres2
--

CREATE TABLE public.orders (
    id integer DEFAULT nextval('public.orders_id_seq'::regclass) NOT NULL,
    user_id character varying(255),
    invoice_id character varying(255),
    created_at timestamp without time zone,
  	paid bool,
  	amount integer
);


ALTER TABLE public.orders OWNER TO postgres2;

--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres2
--

SELECT pg_catalog.setval('public.orders_id_seq', 1, true);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres2
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


INSERT INTO "public"."orders"("user_id","created_at","invoice_id","paid","amount")
VALUES
(1,E'2022-03-14 00:00:00',1,true,1000);
















